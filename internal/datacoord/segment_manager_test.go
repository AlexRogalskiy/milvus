// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datacoord

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	"github.com/milvus-io/milvus-proto/go-api/v2/schemapb"
	etcdkv "github.com/milvus-io/milvus/internal/kv/etcd"
	mockkv "github.com/milvus-io/milvus/internal/kv/mocks"
	"github.com/milvus-io/milvus/internal/metastore/kv/datacoord"
	"github.com/milvus-io/milvus/internal/proto/datapb"
	"github.com/milvus-io/milvus/pkg/util/etcd"
	"github.com/milvus-io/milvus/pkg/util/metautil"
)

func TestManagerOptions(t *testing.T) {
	//	ctx := context.Background()
	Params.Init()
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)
	segmentManager, _ := newSegmentManager(meta, mockAllocator)
	t.Run("test with alloc helper", func(t *testing.T) {
		opt := withAllocHelper(allocHelper{})
		opt.apply(segmentManager)

		assert.True(t, segmentManager.helper.afterCreateSegment == nil)
	})

	t.Run("test withCalUpperLimitPolicy", func(t *testing.T) {
		opt := withCalUpperLimitPolicy(defaultCalUpperLimitPolicy())
		assert.NotNil(t, opt)

		//manual set nil``
		segmentManager.estimatePolicy = nil
		opt.apply(segmentManager)
		assert.True(t, segmentManager.estimatePolicy != nil)
	})

	t.Run("test withAllocPolicy", func(t *testing.T) {
		opt := withAllocPolicy(defaultAllocatePolicy())
		assert.NotNil(t, opt)
		// manual set nil
		segmentManager.allocPolicy = nil
		opt.apply(segmentManager)
		assert.True(t, segmentManager.allocPolicy != nil)
	})

	t.Run("test withSegmentSealPolicy", func(t *testing.T) {
		opt := withSegmentSealPolices(defaultSegmentSealPolicy()...)
		assert.NotNil(t, opt)
		// manual set nil
		segmentManager.segmentSealPolicies = []segmentSealPolicy{}
		opt.apply(segmentManager)
		assert.True(t, len(segmentManager.segmentSealPolicies) > 0)
	})

	t.Run("test withChannelSealPolicies", func(t *testing.T) {
		opt := withChannelSealPolices(getChannelOpenSegCapacityPolicy(1000))
		assert.NotNil(t, opt)
		// manual set nil
		segmentManager.channelSealPolicies = []channelSealPolicy{}
		opt.apply(segmentManager)
		assert.True(t, len(segmentManager.channelSealPolicies) > 0)
	})
	t.Run("test withFlushPolicy", func(t *testing.T) {
		opt := withFlushPolicy(defaultFlushPolicy())
		assert.NotNil(t, opt)
		// manual set nil
		segmentManager.flushPolicy = nil
		opt.apply(segmentManager)
		assert.True(t, segmentManager.flushPolicy != nil)
	})
}

func TestAllocSegment(t *testing.T) {
	ctx := context.Background()
	Params.Init()
	Params.Save(Params.DataCoordCfg.AllocLatestExpireAttempt.Key, "1")
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)
	segmentManager, _ := newSegmentManager(meta, mockAllocator)

	schema := newTestSchema()
	collID, err := mockAllocator.allocID(ctx)
	assert.NoError(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})

	t.Run("normal allocation", func(t *testing.T) {
		allocations, err := segmentManager.AllocSegment(ctx, collID, 100, "c1", 100)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocations))
		assert.EqualValues(t, 100, allocations[0].NumOfRows)
		assert.NotEqualValues(t, 0, allocations[0].SegmentID)
		assert.NotEqualValues(t, 0, allocations[0].ExpireTime)
	})

	t.Run("allocation fails 1", func(t *testing.T) {
		failsAllocator := &FailsAllocator{
			allocTsSucceed: true,
			allocIDSucceed: false,
		}
		segmentManager, err := newSegmentManager(meta, failsAllocator)
		assert.NoError(t, err)
		_, err = segmentManager.AllocSegment(ctx, collID, 100, "c2", 100)
		assert.Error(t, err)
	})

	t.Run("allocation fails 2", func(t *testing.T) {
		failsAllocator := &FailsAllocator{
			allocTsSucceed: false,
			allocIDSucceed: true,
		}
		segmentManager, err := newSegmentManager(meta, failsAllocator)
		assert.Error(t, err)
		assert.Nil(t, segmentManager)
	})
}

func TestLastExpireReset(t *testing.T) {
	//set up meta on dc
	ctx := context.Background()
	Params.Init()
	Params.Save(Params.DataCoordCfg.AllocLatestExpireAttempt.Key, "1")
	Params.Save(Params.DataCoordCfg.SegmentMaxSize.Key, "1")
	mockAllocator := newRootCoordAllocator(newMockRootCoordService())
	etcdCli, _ := etcd.GetEtcdClient(
		Params.EtcdCfg.UseEmbedEtcd.GetAsBool(),
		Params.EtcdCfg.EtcdUseSSL.GetAsBool(),
		Params.EtcdCfg.Endpoints.GetAsStrings(),
		Params.EtcdCfg.EtcdTLSCert.GetValue(),
		Params.EtcdCfg.EtcdTLSKey.GetValue(),
		Params.EtcdCfg.EtcdTLSCACert.GetValue(),
		Params.EtcdCfg.EtcdTLSMinVersion.GetValue())
	rootPath := "/test/segment/last/expire"
	metaKV := etcdkv.NewEtcdKV(etcdCli, rootPath)
	metaKV.RemoveWithPrefix("")
	catalog := datacoord.NewCatalog(metaKV, "", "")
	meta, err := newMeta(context.TODO(), catalog, nil)
	assert.Nil(t, err)
	// add collection
	channelName := "c1"
	schema := newTestSchema()
	collID, err := mockAllocator.allocID(ctx)
	assert.Nil(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})

	//assign segments, set max segment to only 1MB, equalling to 10485 rows
	var bigRows, smallRows int64 = 10000, 1000
	segmentManager, _ := newSegmentManager(meta, mockAllocator)
	allocs, _ := segmentManager.AllocSegment(context.Background(), collID, 0, channelName, bigRows)
	segmentID1, expire1 := allocs[0].SegmentID, allocs[0].ExpireTime
	time.Sleep(100 * time.Millisecond)
	allocs, _ = segmentManager.AllocSegment(context.Background(), collID, 0, channelName, bigRows)
	segmentID2, expire2 := allocs[0].SegmentID, allocs[0].ExpireTime
	time.Sleep(100 * time.Millisecond)
	allocs, _ = segmentManager.AllocSegment(context.Background(), collID, 0, channelName, smallRows)
	segmentID3, expire3 := allocs[0].SegmentID, allocs[0].ExpireTime

	//simulate handleTimeTick op on dataCoord
	meta.SetCurrentRows(segmentID1, bigRows)
	meta.SetCurrentRows(segmentID2, bigRows)
	meta.SetCurrentRows(segmentID3, smallRows)
	segmentManager.tryToSealSegment(expire1, channelName)
	assert.Equal(t, commonpb.SegmentState_Sealed, meta.GetSegment(segmentID1).GetState())
	assert.Equal(t, commonpb.SegmentState_Sealed, meta.GetSegment(segmentID2).GetState())
	assert.Equal(t, commonpb.SegmentState_Growing, meta.GetSegment(segmentID3).GetState())

	//pretend that dataCoord break down
	metaKV.Close()
	etcdCli.Close()

	//dataCoord restart
	newEtcdCli, _ := etcd.GetEtcdClient(Params.EtcdCfg.UseEmbedEtcd.GetAsBool(), Params.EtcdCfg.EtcdUseSSL.GetAsBool(),
		Params.EtcdCfg.Endpoints.GetAsStrings(), Params.EtcdCfg.EtcdTLSCert.GetValue(),
		Params.EtcdCfg.EtcdTLSKey.GetValue(), Params.EtcdCfg.EtcdTLSCACert.GetValue(), Params.EtcdCfg.EtcdTLSMinVersion.GetValue())
	newMetaKV := etcdkv.NewEtcdKV(newEtcdCli, rootPath)
	defer newMetaKV.RemoveWithPrefix("")
	newCatalog := datacoord.NewCatalog(newMetaKV, "", "")
	restartedMeta, err := newMeta(context.TODO(), newCatalog, nil)
	restartedMeta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
	assert.Nil(t, err)
	newSegmentManager, _ := newSegmentManager(restartedMeta, mockAllocator)
	//reset row number to avoid being cleaned by empty segment
	restartedMeta.SetCurrentRows(segmentID1, bigRows)
	restartedMeta.SetCurrentRows(segmentID2, bigRows)
	restartedMeta.SetCurrentRows(segmentID3, smallRows)

	//verify lastExpire of growing and sealed segments
	segment1, segment2, segment3 := restartedMeta.GetSegment(segmentID1), restartedMeta.GetSegment(segmentID2), restartedMeta.GetSegment(segmentID3)
	//segmentState should not be altered but growing segment's lastExpire has been reset to the latest
	assert.Equal(t, commonpb.SegmentState_Sealed, segment1.GetState())
	assert.Equal(t, commonpb.SegmentState_Sealed, segment2.GetState())
	assert.Equal(t, commonpb.SegmentState_Growing, segment3.GetState())
	assert.Equal(t, expire1, segment1.GetLastExpireTime())
	assert.Equal(t, expire2, segment2.GetLastExpireTime())
	assert.True(t, segment3.GetLastExpireTime() > expire3)
	flushableSegIds, _ := newSegmentManager.GetFlushableSegments(context.Background(), channelName, expire3)
	assert.ElementsMatch(t, []UniqueID{segmentID1, segmentID2}, flushableSegIds) // segment1 and segment2 can be flushed
	newAlloc, err := newSegmentManager.AllocSegment(context.Background(), collID, 0, channelName, 2000)
	assert.Nil(t, err)
	assert.Equal(t, segmentID3, newAlloc[0].SegmentID) // segment3 still can be used to allocate
}

func TestAllocSegmentForImport(t *testing.T) {
	ctx := context.Background()
	Params.Init()
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)
	segmentManager, _ := newSegmentManager(meta, mockAllocator)

	schema := newTestSchema()
	collID, err := mockAllocator.allocID(ctx)
	assert.NoError(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})

	t.Run("normal allocation", func(t *testing.T) {
		allocation, err := segmentManager.allocSegmentForImport(ctx, collID, 100, "c1", 100, 0)
		assert.NoError(t, err)
		assert.NotNil(t, allocation)
		assert.EqualValues(t, 100, allocation.NumOfRows)
		assert.NotEqualValues(t, 0, allocation.SegmentID)
		assert.NotEqualValues(t, 0, allocation.ExpireTime)
	})

	t.Run("allocation fails 1", func(t *testing.T) {
		failsAllocator := &FailsAllocator{
			allocTsSucceed: true,
		}
		segmentManager, _ := newSegmentManager(meta, failsAllocator)
		_, err := segmentManager.allocSegmentForImport(ctx, collID, 100, "c1", 100, 0)
		assert.Error(t, err)
	})

	t.Run("allocation fails 2", func(t *testing.T) {
		failsAllocator := &FailsAllocator{
			allocIDSucceed: true,
		}
		segmentManager, _ := newSegmentManager(meta, failsAllocator)
		_, err := segmentManager.allocSegmentForImport(ctx, collID, 100, "c1", 100, 0)
		assert.Error(t, err)
	})
}

func TestLoadSegmentsFromMeta(t *testing.T) {
	ctx := context.Background()
	Params.Init()
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)

	schema := newTestSchema()
	collID, err := mockAllocator.allocID(ctx)
	assert.NoError(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})

	sealedSegment := &datapb.SegmentInfo{
		ID:             1,
		CollectionID:   collID,
		PartitionID:    0,
		InsertChannel:  "",
		State:          commonpb.SegmentState_Sealed,
		MaxRowNum:      100,
		LastExpireTime: 1000,
	}
	growingSegment := &datapb.SegmentInfo{
		ID:             2,
		CollectionID:   collID,
		PartitionID:    0,
		InsertChannel:  "",
		State:          commonpb.SegmentState_Growing,
		MaxRowNum:      100,
		LastExpireTime: 1000,
	}
	flushedSegment := &datapb.SegmentInfo{
		ID:             3,
		CollectionID:   collID,
		PartitionID:    0,
		InsertChannel:  "",
		State:          commonpb.SegmentState_Flushed,
		MaxRowNum:      100,
		LastExpireTime: 1000,
	}
	err = meta.AddSegment(NewSegmentInfo(sealedSegment))
	assert.NoError(t, err)
	err = meta.AddSegment(NewSegmentInfo(growingSegment))
	assert.NoError(t, err)
	err = meta.AddSegment(NewSegmentInfo(flushedSegment))
	assert.NoError(t, err)

	segmentManager, _ := newSegmentManager(meta, mockAllocator)
	segments := segmentManager.segments
	assert.EqualValues(t, 2, len(segments))
}

func TestSaveSegmentsToMeta(t *testing.T) {
	Params.Init()
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)

	schema := newTestSchema()
	collID, err := mockAllocator.allocID(context.Background())
	assert.NoError(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
	segmentManager, _ := newSegmentManager(meta, mockAllocator)
	allocations, err := segmentManager.AllocSegment(context.Background(), collID, 0, "c1", 1000)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(allocations))
	_, err = segmentManager.SealAllSegments(context.Background(), collID, nil, false)
	assert.NoError(t, err)
	segment := meta.GetHealthySegment(allocations[0].SegmentID)
	assert.NotNil(t, segment)
	assert.EqualValues(t, segment.LastExpireTime, allocations[0].ExpireTime)
	assert.EqualValues(t, commonpb.SegmentState_Sealed, segment.State)
}

func TestSaveSegmentsToMetaWithSpecificSegments(t *testing.T) {
	Params.Init()
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)

	schema := newTestSchema()
	collID, err := mockAllocator.allocID(context.Background())
	assert.NoError(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
	segmentManager, _ := newSegmentManager(meta, mockAllocator)
	allocations, err := segmentManager.AllocSegment(context.Background(), collID, 0, "c1", 1000)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(allocations))
	_, err = segmentManager.SealAllSegments(context.Background(), collID, []int64{allocations[0].SegmentID}, false)
	assert.NoError(t, err)
	segment := meta.GetHealthySegment(allocations[0].SegmentID)
	assert.NotNil(t, segment)
	assert.EqualValues(t, segment.LastExpireTime, allocations[0].ExpireTime)
	assert.EqualValues(t, commonpb.SegmentState_Sealed, segment.State)
}

func TestDropSegment(t *testing.T) {
	Params.Init()
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)

	schema := newTestSchema()
	collID, err := mockAllocator.allocID(context.Background())
	assert.NoError(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
	segmentManager, _ := newSegmentManager(meta, mockAllocator)
	allocations, err := segmentManager.AllocSegment(context.Background(), collID, 0, "c1", 1000)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(allocations))
	segID := allocations[0].SegmentID
	segment := meta.GetHealthySegment(segID)
	assert.NotNil(t, segment)

	segmentManager.DropSegment(context.Background(), segID)
	segment = meta.GetHealthySegment(segID)
	assert.NotNil(t, segment)
}

func TestAllocRowsLargerThanOneSegment(t *testing.T) {
	Params.Init()
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)

	schema := newTestSchema()
	collID, err := mockAllocator.allocID(context.Background())
	assert.NoError(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})

	var mockPolicy = func(schema *schemapb.CollectionSchema) (int, error) {
		return 1, nil
	}
	segmentManager, _ := newSegmentManager(meta, mockAllocator, withCalUpperLimitPolicy(mockPolicy))
	allocations, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "c1", 2)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(allocations))
	assert.EqualValues(t, 1, allocations[0].NumOfRows)
	assert.EqualValues(t, 1, allocations[1].NumOfRows)
}

func TestExpireAllocation(t *testing.T) {
	Params.Init()
	mockAllocator := newMockAllocator()
	meta, err := newMemoryMeta()
	assert.NoError(t, err)

	schema := newTestSchema()
	collID, err := mockAllocator.allocID(context.Background())
	assert.NoError(t, err)
	meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})

	var mockPolicy = func(schema *schemapb.CollectionSchema) (int, error) {
		return 10000000, nil
	}
	segmentManager, _ := newSegmentManager(meta, mockAllocator, withCalUpperLimitPolicy(mockPolicy))
	// alloc 100 times and expire
	var maxts Timestamp
	var id int64 = -1
	for i := 0; i < 100; i++ {
		allocs, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "ch1", 100)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocs))
		if id == -1 {
			id = allocs[0].SegmentID
		} else {
			assert.EqualValues(t, id, allocs[0].SegmentID)
		}
		if allocs[0].ExpireTime > maxts {
			maxts = allocs[0].ExpireTime
		}
	}

	segment := meta.GetHealthySegment(id)
	assert.NotNil(t, segment)
	assert.EqualValues(t, 100, len(segment.allocations))
	err = segmentManager.ExpireAllocations("ch1", maxts)
	assert.NoError(t, err)
	segment = meta.GetHealthySegment(id)
	assert.NotNil(t, segment)
	assert.EqualValues(t, 0, len(segment.allocations))
}

func TestCleanExpiredBulkloadSegment(t *testing.T) {
	t.Run("expiredBulkloadSegment", func(t *testing.T) {
		Params.Init()
		mockAllocator := newMockAllocator()
		meta, err := newMemoryMeta()
		assert.NoError(t, err)

		schema := newTestSchema()
		collID, err := mockAllocator.allocID(context.Background())
		assert.NoError(t, err)
		meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
		segmentManager, _ := newSegmentManager(meta, mockAllocator)
		allocation, err := segmentManager.allocSegmentForImport(context.TODO(), collID, 0, "c1", 2, 1)
		assert.NoError(t, err)

		ids, err := segmentManager.GetFlushableSegments(context.TODO(), "c1", allocation.ExpireTime)
		assert.NoError(t, err)
		assert.EqualValues(t, len(ids), 0)

		assert.EqualValues(t, len(segmentManager.segments), 1)

		ids, err = segmentManager.GetFlushableSegments(context.TODO(), "c1", allocation.ExpireTime+1)
		assert.NoError(t, err)
		assert.Empty(t, ids)
		assert.EqualValues(t, len(ids), 0)

		assert.EqualValues(t, len(segmentManager.segments), 0)
	})
}

func TestGetFlushableSegments(t *testing.T) {
	t.Run("get flushable segments between small interval", func(t *testing.T) {
		Params.Init()
		mockAllocator := newMockAllocator()
		meta, err := newMemoryMeta()
		assert.NoError(t, err)

		schema := newTestSchema()
		collID, err := mockAllocator.allocID(context.Background())
		assert.NoError(t, err)
		meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
		segmentManager, _ := newSegmentManager(meta, mockAllocator)
		allocations, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "c1", 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocations))

		ids, err := segmentManager.SealAllSegments(context.TODO(), collID, nil, false)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(ids))
		assert.EqualValues(t, allocations[0].SegmentID, ids[0])

		meta.SetCurrentRows(allocations[0].SegmentID, 1)
		ids, err = segmentManager.GetFlushableSegments(context.TODO(), "c1", allocations[0].ExpireTime)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(ids))
		assert.EqualValues(t, allocations[0].SegmentID, ids[0])

		meta.SetLastFlushTime(allocations[0].SegmentID, time.Now())
		ids, err = segmentManager.GetFlushableSegments(context.TODO(), "c1", allocations[0].ExpireTime)
		assert.NoError(t, err)
		assert.Empty(t, ids)

		meta.SetLastFlushTime(allocations[0].SegmentID, time.Now().Local().Add(-flushInterval))
		ids, err = segmentManager.GetFlushableSegments(context.TODO(), "c1", allocations[0].ExpireTime)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(ids))
		assert.EqualValues(t, allocations[0].SegmentID, ids[0])

		meta.SetCurrentRows(allocations[0].SegmentID, 0)
		ids, err = segmentManager.GetFlushableSegments(context.TODO(), "c1", allocations[0].ExpireTime)
		assert.NoError(t, err)
		assert.Empty(t, ids)
		assert.Nil(t, meta.GetHealthySegment(allocations[0].SegmentID))
	})
}

func TestTryToSealSegment(t *testing.T) {
	t.Run("normal seal with segment policies", func(t *testing.T) {
		Params.Init()
		mockAllocator := newMockAllocator()
		meta, err := newMemoryMeta()
		assert.NoError(t, err)

		schema := newTestSchema()
		collID, err := mockAllocator.allocID(context.Background())
		assert.NoError(t, err)
		meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
		segmentManager, _ := newSegmentManager(meta, mockAllocator, withSegmentSealPolices(sealByLifetimePolicy(math.MinInt64))) //always seal
		allocations, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "c1", 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocations))

		ts, err := segmentManager.allocator.allocTimestamp(context.Background())
		assert.NoError(t, err)
		err = segmentManager.tryToSealSegment(ts, "c1")
		assert.NoError(t, err)

		for _, seg := range segmentManager.meta.segments.segments {
			assert.Equal(t, commonpb.SegmentState_Sealed, seg.GetState())
		}
	})

	t.Run("normal seal with channel seal policies", func(t *testing.T) {
		Params.Init()
		mockAllocator := newMockAllocator()
		meta, err := newMemoryMeta()
		assert.NoError(t, err)

		schema := newTestSchema()
		collID, err := mockAllocator.allocID(context.Background())
		assert.NoError(t, err)
		meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
		segmentManager, _ := newSegmentManager(meta, mockAllocator, withChannelSealPolices(getChannelOpenSegCapacityPolicy(-1))) //always seal
		allocations, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "c1", 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocations))

		ts, err := segmentManager.allocator.allocTimestamp(context.Background())
		assert.NoError(t, err)
		err = segmentManager.tryToSealSegment(ts, "c1")
		assert.NoError(t, err)

		for _, seg := range segmentManager.meta.segments.segments {
			assert.Equal(t, commonpb.SegmentState_Sealed, seg.GetState())
		}
	})

	t.Run("normal seal with both segment & channel seal policy", func(t *testing.T) {
		Params.Init()
		mockAllocator := newMockAllocator()
		meta, err := newMemoryMeta()
		assert.NoError(t, err)

		schema := newTestSchema()
		collID, err := mockAllocator.allocID(context.Background())
		assert.NoError(t, err)
		meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
		segmentManager, _ := newSegmentManager(meta, mockAllocator,
			withSegmentSealPolices(sealByLifetimePolicy(math.MinInt64)),
			withChannelSealPolices(getChannelOpenSegCapacityPolicy(-1))) //always seal
		allocations, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "c1", 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocations))

		ts, err := segmentManager.allocator.allocTimestamp(context.Background())
		assert.NoError(t, err)
		err = segmentManager.tryToSealSegment(ts, "c1")
		assert.NoError(t, err)

		for _, seg := range segmentManager.meta.segments.segments {
			assert.Equal(t, commonpb.SegmentState_Sealed, seg.GetState())
		}
	})

	t.Run("test sealByMaxBinlogFileNumberPolicy", func(t *testing.T) {
		Params.Init()
		mockAllocator := newMockAllocator()
		meta, err := newMemoryMeta()
		assert.NoError(t, err)

		schema := newTestSchema()
		collID, err := mockAllocator.allocID(context.Background())
		assert.NoError(t, err)
		meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
		segmentManager, _ := newSegmentManager(meta, mockAllocator)
		allocations, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "c1", 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocations))

		ts, err := segmentManager.allocator.allocTimestamp(context.Background())
		assert.NoError(t, err)

		// No seal polices
		{
			err = segmentManager.tryToSealSegment(ts, "c1")
			assert.NoError(t, err)
			segments := segmentManager.meta.segments.segments
			assert.Equal(t, 1, len(segments))
			for _, seg := range segments {
				assert.Equal(t, commonpb.SegmentState_Growing, seg.GetState())
			}
		}

		// Not trigger seal
		{
			segmentManager.segmentSealPolicies = []segmentSealPolicy{sealByMaxBinlogFileNumberPolicy(2)}
			segments := segmentManager.meta.segments.segments
			assert.Equal(t, 1, len(segments))
			for _, seg := range segments {
				seg.Statslogs = []*datapb.FieldBinlog{
					{
						FieldID: 2,
						Binlogs: []*datapb.Binlog{
							{
								EntriesNum: 10,
								LogID:      3,
								LogPath:    metautil.BuildInsertLogPath("", 1, 1, seg.ID, 2, 3),
							},
						},
					},
				}
				err = segmentManager.tryToSealSegment(ts, "c1")
				assert.NoError(t, err)
				seg = segmentManager.meta.segments.segments[seg.ID]
				assert.Equal(t, commonpb.SegmentState_Growing, seg.GetState())
			}
		}

		// Trigger seal
		{
			segmentManager.segmentSealPolicies = []segmentSealPolicy{sealByMaxBinlogFileNumberPolicy(2)}
			segments := segmentManager.meta.segments.segments
			assert.Equal(t, 1, len(segments))
			for _, seg := range segments {
				seg.Statslogs = []*datapb.FieldBinlog{
					{
						FieldID: 1,
						Binlogs: []*datapb.Binlog{
							{
								EntriesNum: 10,
								LogID:      1,
								LogPath:    metautil.BuildInsertLogPath("", 1, 1, seg.ID, 1, 3),
							},
							{
								EntriesNum: 20,
								LogID:      2,
								LogPath:    metautil.BuildInsertLogPath("", 1, 1, seg.ID, 1, 2),
							},
						},
					},
				}
				err = segmentManager.tryToSealSegment(ts, "c1")
				assert.NoError(t, err)
				seg = segmentManager.meta.segments.segments[seg.ID]
				assert.Equal(t, commonpb.SegmentState_Sealed, seg.GetState())
			}
		}
	})

	t.Run("seal with segment policy with kv fails", func(t *testing.T) {
		Params.Init()
		mockAllocator := newMockAllocator()
		memoryKV := NewMetaMemoryKV()
		catalog := datacoord.NewCatalog(memoryKV, "", "")
		meta, err := newMeta(context.TODO(), catalog, nil)
		assert.NoError(t, err)

		schema := newTestSchema()
		collID, err := mockAllocator.allocID(context.Background())
		assert.NoError(t, err)
		meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
		segmentManager, _ := newSegmentManager(meta, mockAllocator, withSegmentSealPolices(sealByLifetimePolicy(math.MinInt64))) //always seal
		allocations, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "c1", 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocations))

		metakv := mockkv.NewMetaKv(t)
		metakv.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("failed")).Maybe()
		metakv.EXPECT().MultiSave(mock.Anything).Return(errors.New("failed")).Maybe()
		metakv.EXPECT().LoadWithPrefix(mock.Anything).Return(nil, nil, nil).Maybe()
		segmentManager.meta.catalog = &datacoord.Catalog{MetaKv: metakv}

		ts, err := segmentManager.allocator.allocTimestamp(context.Background())
		assert.NoError(t, err)
		err = segmentManager.tryToSealSegment(ts, "c1")
		assert.Error(t, err)
	})

	t.Run("seal with channel policy with kv fails", func(t *testing.T) {
		Params.Init()
		mockAllocator := newMockAllocator()
		memoryKV := NewMetaMemoryKV()
		catalog := datacoord.NewCatalog(memoryKV, "", "")
		meta, err := newMeta(context.TODO(), catalog, nil)
		assert.NoError(t, err)

		schema := newTestSchema()
		collID, err := mockAllocator.allocID(context.Background())
		assert.NoError(t, err)
		meta.AddCollection(&collectionInfo{ID: collID, Schema: schema})
		segmentManager, _ := newSegmentManager(meta, mockAllocator, withChannelSealPolices(getChannelOpenSegCapacityPolicy(-1))) //always seal
		allocations, err := segmentManager.AllocSegment(context.TODO(), collID, 0, "c1", 2)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(allocations))

		metakv := mockkv.NewMetaKv(t)
		metakv.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("failed")).Maybe()
		metakv.EXPECT().MultiSave(mock.Anything).Return(errors.New("failed")).Maybe()
		metakv.EXPECT().LoadWithPrefix(mock.Anything).Return(nil, nil, nil).Maybe()
		segmentManager.meta.catalog = &datacoord.Catalog{MetaKv: metakv}

		ts, err := segmentManager.allocator.allocTimestamp(context.Background())
		assert.NoError(t, err)
		err = segmentManager.tryToSealSegment(ts, "c1")
		assert.Error(t, err)
	})
}

func TestAllocationPool(t *testing.T) {
	t.Run("normal get&put", func(t *testing.T) {
		allocPool = sync.Pool{
			New: func() interface{} {
				return &Allocation{}
			},
		}

		allo := getAllocation(100)
		assert.EqualValues(t, 100, allo.NumOfRows)
		assert.EqualValues(t, 0, allo.ExpireTime)
		assert.EqualValues(t, 0, allo.SegmentID)

		putAllocation(allo)
	})

	t.Run("put nil", func(t *testing.T) {
		var allo *Allocation
		allocPool = sync.Pool{
			New: func() interface{} {
				return &Allocation{}
			},
		}
		putAllocation(allo)
		allo = getAllocation(100)
		assert.EqualValues(t, 100, allo.NumOfRows)
		assert.EqualValues(t, 0, allo.ExpireTime)
		assert.EqualValues(t, 0, allo.SegmentID)
	})

	t.Run("put something else", func(t *testing.T) {
		allocPool = sync.Pool{
			New: func() interface{} {
				return &Allocation{}
			},
		}
		allocPool.Put(&struct{}{})
		allo := getAllocation(100)
		assert.EqualValues(t, 100, allo.NumOfRows)
		assert.EqualValues(t, 0, allo.ExpireTime)
		assert.EqualValues(t, 0, allo.SegmentID)

	})
}

func TestSegmentManager_DropSegmentsOfChannel(t *testing.T) {
	type fields struct {
		meta     *meta
		segments []UniqueID
	}
	type args struct {
		channel string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []UniqueID
	}{
		{
			"test drop segments",
			fields{
				meta: &meta{
					segments: &SegmentsInfo{
						segments: map[int64]*SegmentInfo{
							1: {
								SegmentInfo: &datapb.SegmentInfo{
									ID:            1,
									InsertChannel: "ch1",
									State:         commonpb.SegmentState_Flushed,
								},
							},
							2: {
								SegmentInfo: &datapb.SegmentInfo{
									ID:            2,
									InsertChannel: "ch2",
									State:         commonpb.SegmentState_Flushed,
								},
							},
						},
					},
				},
				segments: []UniqueID{1, 2},
			},
			args{
				"ch1",
			},
			[]UniqueID{2},
		},
		{
			"test drop segments with dropped segment",
			fields{
				meta: &meta{
					segments: &SegmentsInfo{
						segments: map[int64]*SegmentInfo{
							1: {
								SegmentInfo: &datapb.SegmentInfo{
									ID:            1,
									InsertChannel: "ch1",
									State:         commonpb.SegmentState_Dropped,
								},
							},
							2: {
								SegmentInfo: &datapb.SegmentInfo{
									ID:            2,
									InsertChannel: "ch2",
									State:         commonpb.SegmentState_Growing,
								},
							},
						},
					},
				},
				segments: []UniqueID{1, 2, 3},
			},
			args{
				"ch1",
			},
			[]UniqueID{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SegmentManager{
				meta:     tt.fields.meta,
				segments: tt.fields.segments,
			}
			s.DropSegmentsOfChannel(context.TODO(), tt.args.channel)
			assert.ElementsMatch(t, tt.want, s.segments)
		})
	}
}
