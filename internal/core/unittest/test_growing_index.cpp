// Copyright (C) 2019-2020 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License

#include <gtest/gtest.h>

#include "pb/plan.pb.h"
#include "segcore/SegmentGrowing.h"
#include "segcore/SegmentGrowingImpl.h"
#include "pb/schema.pb.h"
#include "test_utils/DataGen.h"
#include "query/Plan.h"

using namespace milvus::segcore;
using namespace milvus;
namespace pb = milvus::proto;

TEST(GrowingIndex, Correctness) {
    auto schema = std::make_shared<Schema>();
    auto pk = schema->AddDebugField("pk", DataType::INT64);
    auto random = schema->AddDebugField("random", DataType::DOUBLE);
    auto vec = schema->AddDebugField(
        "embeddings", DataType::VECTOR_FLOAT, 128, knowhere::metric::L2);
    schema->set_primary_field_id(pk);

    std::map<std::string, std::string> index_params = {
        {"index_type", "IVF_FLAT"}, {"metric_type", "L2"}, {"nlist", "128"}};
    std::map<std::string, std::string> type_params = {{"dim", "128"}};
    FieldIndexMeta fieldIndexMeta(
        vec, std::move(index_params), std::move(type_params));
    auto& config = SegcoreConfig::default_config();
    config.set_chunk_rows(1024);
    config.set_enable_growing_segment_index(true);
    std::map<FieldId, FieldIndexMeta> filedMap = {{vec, fieldIndexMeta}};
    IndexMetaPtr metaPtr =
        std::make_shared<CollectionIndexMeta>(226985, std::move(filedMap));
    auto segment = CreateGrowingSegment(schema, metaPtr);
    auto segmentImplPtr = dynamic_cast<SegmentGrowingImpl*>(segment.get());

    milvus::proto::plan::PlanNode plan_node;
    auto vector_anns = plan_node.mutable_vector_anns();
    vector_anns->set_is_binary(false);
    vector_anns->set_placeholder_tag("$0");
    vector_anns->set_field_id(102);
    auto query_info = vector_anns->mutable_query_info();
    query_info->set_topk(5);
    query_info->set_round_decimal(3);
    query_info->set_metric_type("l2");
    query_info->set_search_params(R"({"nprobe": 16})");
    auto plan_str = plan_node.SerializeAsString();

    milvus::proto::plan::PlanNode range_query_plan_node;
    auto vector_range_querys = range_query_plan_node.mutable_vector_anns();
    vector_range_querys->set_is_binary(false);
    vector_range_querys->set_placeholder_tag("$0");
    vector_range_querys->set_field_id(102);
    auto range_query_info = vector_range_querys->mutable_query_info();
    range_query_info->set_topk(5);
    range_query_info->set_round_decimal(3);
    range_query_info->set_metric_type("l2");
    range_query_info->set_search_params(
        R"({"nprobe": 10, "radius": 600, "range_filter": 500})");
    auto range_plan_str = range_query_plan_node.SerializeAsString();

    int64_t per_batch = 10000;
    int64_t n_batch = 20;
    int64_t top_k = 5;
    for (int64_t i = 0; i < n_batch; i++) {
        auto dataset = DataGen(schema, per_batch);
        auto offset = segment->PreInsert(per_batch);
        auto pks = dataset.get_col<int64_t>(pk);
        segment->Insert(offset,
                        per_batch,
                        dataset.row_ids_.data(),
                        dataset.timestamps_.data(),
                        dataset.raw_);
        auto filed_data = segmentImplPtr->get_insert_record()
                              .get_field_data<milvus::FloatVector>(vec);

        auto inserted = (i + 1) * per_batch;
        //once index built, chunk data will be removed
        if (i < 2) {
            EXPECT_EQ(filed_data->num_chunk(),
                      upper_div(inserted, filed_data->get_size_per_chunk()));
        } else {
            EXPECT_EQ(filed_data->num_chunk(), 0);
        }

        auto num_queries = 5;
        auto ph_group_raw = CreatePlaceholderGroup(num_queries, 128, 1024);

        auto plan = milvus::query::CreateSearchPlanByExpr(
            *schema, plan_str.data(), plan_str.size());
        auto ph_group =
            ParsePlaceholderGroup(plan.get(), ph_group_raw.SerializeAsString());
        auto sr = segment->Search(plan.get(), ph_group.get());
        EXPECT_EQ(sr->total_nq_, num_queries);
        EXPECT_EQ(sr->unity_topK_, top_k);
        EXPECT_EQ(sr->distances_.size(), num_queries * top_k);
        EXPECT_EQ(sr->seg_offsets_.size(), num_queries * top_k);

        auto range_plan = milvus::query::CreateSearchPlanByExpr(
            *schema, range_plan_str.data(), range_plan_str.size());
        auto range_ph_group = ParsePlaceholderGroup(
            range_plan.get(), ph_group_raw.SerializeAsString());
        auto range_sr = segment->Search(range_plan.get(), range_ph_group.get());
        ASSERT_EQ(range_sr->total_nq_, num_queries);
        EXPECT_EQ(sr->unity_topK_, top_k);
        EXPECT_EQ(sr->distances_.size(), num_queries * top_k);
        EXPECT_EQ(sr->seg_offsets_.size(), num_queries * top_k);
        for (int j = 0; j < range_sr->seg_offsets_.size(); j++) {
            if (range_sr->seg_offsets_[j] != -1) {
                EXPECT_TRUE(sr->distances_[j] >= 500.0 &&
                            sr->distances_[j] <= 600.0);
            }
        }
    }
}

using Param = const char*;

class GrowingIndexGetVectorTest : public ::testing::TestWithParam<Param> {
    void
    SetUp() override {
        auto param = GetParam();
        metricType = param;
    }

 protected:
    const char* metricType;
};

INSTANTIATE_TEST_CASE_P(IndexTypeParameters,
                        GrowingIndexGetVectorTest,
                        ::testing::Values(knowhere::metric::L2,
                                          knowhere::metric::COSINE,
                                          knowhere::metric::IP));

TEST_P(GrowingIndexGetVectorTest, GetVector) {
    auto schema = std::make_shared<Schema>();
    auto pk = schema->AddDebugField("pk", DataType::INT64);
    auto random = schema->AddDebugField("random", DataType::DOUBLE);
    auto vec = schema->AddDebugField(
        "embeddings", DataType::VECTOR_FLOAT, 128, metricType);
    schema->set_primary_field_id(pk);

    std::map<std::string, std::string> index_params = {
        {"index_type", "IVF_FLAT"},
        {"metric_type", metricType},
        {"nlist", "128"}};
    std::map<std::string, std::string> type_params = {{"dim", "128"}};
    FieldIndexMeta fieldIndexMeta(
        vec, std::move(index_params), std::move(type_params));
    auto& config = SegcoreConfig::default_config();
    config.set_chunk_rows(1024);
    config.set_enable_growing_segment_index(true);
    std::map<FieldId, FieldIndexMeta> filedMap = {{vec, fieldIndexMeta}};
    IndexMetaPtr metaPtr =
        std::make_shared<CollectionIndexMeta>(100000, std::move(filedMap));
    auto segment_growing = CreateGrowingSegment(schema, metaPtr);
    auto segment = dynamic_cast<SegmentGrowingImpl*>(segment_growing.get());

    int64_t per_batch = 5000;
    int64_t n_batch = 20;
    int64_t dim = 128;
    for (int64_t i = 0; i < n_batch; i++) {
        auto dataset = DataGen(schema, per_batch);
        auto fakevec = dataset.get_col<float>(vec);
        auto offset = segment->PreInsert(per_batch);
        segment->Insert(offset,
                        per_batch,
                        dataset.row_ids_.data(),
                        dataset.timestamps_.data(),
                        dataset.raw_);
        auto num_inserted = (i + 1) * per_batch;
        auto ids_ds = GenRandomIds(num_inserted);
        auto result =
            segment->bulk_subscript(vec, ids_ds->GetIds(), num_inserted);

        auto vector = result.get()->mutable_vectors()->float_vector().data();
        EXPECT_TRUE(vector.size() == num_inserted * dim);
        for (size_t i = 0; i < num_inserted; ++i) {
            auto id = ids_ds->GetIds()[i];
            for (size_t j = 0; j < 128; ++j) {
                EXPECT_TRUE(vector[i * dim + j] ==
                            fakevec[(id % per_batch) * dim + j]);
            }
        }
    }
}
