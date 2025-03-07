// Code generated by mockery v2.32.4. DO NOT EDIT.

package segments

import (
	commonpb "github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	mock "github.com/stretchr/testify/mock"

	msgpb "github.com/milvus-io/milvus-proto/go-api/v2/msgpb"

	segcorepb "github.com/milvus-io/milvus/internal/proto/segcorepb"

	storage "github.com/milvus-io/milvus/internal/storage"
)

// MockSegment is an autogenerated mock type for the Segment type
type MockSegment struct {
	mock.Mock
}

type MockSegment_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSegment) EXPECT() *MockSegment_Expecter {
	return &MockSegment_Expecter{mock: &_m.Mock}
}

// AddIndex provides a mock function with given fields: fieldID, index
func (_m *MockSegment) AddIndex(fieldID int64, index *IndexedFieldInfo) {
	_m.Called(fieldID, index)
}

// MockSegment_AddIndex_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddIndex'
type MockSegment_AddIndex_Call struct {
	*mock.Call
}

// AddIndex is a helper method to define mock.On call
//   - fieldID int64
//   - index *IndexedFieldInfo
func (_e *MockSegment_Expecter) AddIndex(fieldID interface{}, index interface{}) *MockSegment_AddIndex_Call {
	return &MockSegment_AddIndex_Call{Call: _e.mock.On("AddIndex", fieldID, index)}
}

func (_c *MockSegment_AddIndex_Call) Run(run func(fieldID int64, index *IndexedFieldInfo)) *MockSegment_AddIndex_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(*IndexedFieldInfo))
	})
	return _c
}

func (_c *MockSegment_AddIndex_Call) Return() *MockSegment_AddIndex_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockSegment_AddIndex_Call) RunAndReturn(run func(int64, *IndexedFieldInfo)) *MockSegment_AddIndex_Call {
	_c.Call.Return(run)
	return _c
}

// CASVersion provides a mock function with given fields: _a0, _a1
func (_m *MockSegment) CASVersion(_a0 int64, _a1 int64) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64, int64) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockSegment_CASVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CASVersion'
type MockSegment_CASVersion_Call struct {
	*mock.Call
}

// CASVersion is a helper method to define mock.On call
//   - _a0 int64
//   - _a1 int64
func (_e *MockSegment_Expecter) CASVersion(_a0 interface{}, _a1 interface{}) *MockSegment_CASVersion_Call {
	return &MockSegment_CASVersion_Call{Call: _e.mock.On("CASVersion", _a0, _a1)}
}

func (_c *MockSegment_CASVersion_Call) Run(run func(_a0 int64, _a1 int64)) *MockSegment_CASVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(int64))
	})
	return _c
}

func (_c *MockSegment_CASVersion_Call) Return(_a0 bool) *MockSegment_CASVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_CASVersion_Call) RunAndReturn(run func(int64, int64) bool) *MockSegment_CASVersion_Call {
	_c.Call.Return(run)
	return _c
}

// Collection provides a mock function with given fields:
func (_m *MockSegment) Collection() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockSegment_Collection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Collection'
type MockSegment_Collection_Call struct {
	*mock.Call
}

// Collection is a helper method to define mock.On call
func (_e *MockSegment_Expecter) Collection() *MockSegment_Collection_Call {
	return &MockSegment_Collection_Call{Call: _e.mock.On("Collection")}
}

func (_c *MockSegment_Collection_Call) Run(run func()) *MockSegment_Collection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_Collection_Call) Return(_a0 int64) *MockSegment_Collection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_Collection_Call) RunAndReturn(run func() int64) *MockSegment_Collection_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: primaryKeys, timestamps
func (_m *MockSegment) Delete(primaryKeys []storage.PrimaryKey, timestamps []uint64) error {
	ret := _m.Called(primaryKeys, timestamps)

	var r0 error
	if rf, ok := ret.Get(0).(func([]storage.PrimaryKey, []uint64) error); ok {
		r0 = rf(primaryKeys, timestamps)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSegment_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockSegment_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - primaryKeys []storage.PrimaryKey
//   - timestamps []uint64
func (_e *MockSegment_Expecter) Delete(primaryKeys interface{}, timestamps interface{}) *MockSegment_Delete_Call {
	return &MockSegment_Delete_Call{Call: _e.mock.On("Delete", primaryKeys, timestamps)}
}

func (_c *MockSegment_Delete_Call) Run(run func(primaryKeys []storage.PrimaryKey, timestamps []uint64)) *MockSegment_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]storage.PrimaryKey), args[1].([]uint64))
	})
	return _c
}

func (_c *MockSegment_Delete_Call) Return(_a0 error) *MockSegment_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_Delete_Call) RunAndReturn(run func([]storage.PrimaryKey, []uint64) error) *MockSegment_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// ExistIndex provides a mock function with given fields: fieldID
func (_m *MockSegment) ExistIndex(fieldID int64) bool {
	ret := _m.Called(fieldID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(int64) bool); ok {
		r0 = rf(fieldID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockSegment_ExistIndex_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistIndex'
type MockSegment_ExistIndex_Call struct {
	*mock.Call
}

// ExistIndex is a helper method to define mock.On call
//   - fieldID int64
func (_e *MockSegment_Expecter) ExistIndex(fieldID interface{}) *MockSegment_ExistIndex_Call {
	return &MockSegment_ExistIndex_Call{Call: _e.mock.On("ExistIndex", fieldID)}
}

func (_c *MockSegment_ExistIndex_Call) Run(run func(fieldID int64)) *MockSegment_ExistIndex_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockSegment_ExistIndex_Call) Return(_a0 bool) *MockSegment_ExistIndex_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_ExistIndex_Call) RunAndReturn(run func(int64) bool) *MockSegment_ExistIndex_Call {
	_c.Call.Return(run)
	return _c
}

// GetIndex provides a mock function with given fields: fieldID
func (_m *MockSegment) GetIndex(fieldID int64) *IndexedFieldInfo {
	ret := _m.Called(fieldID)

	var r0 *IndexedFieldInfo
	if rf, ok := ret.Get(0).(func(int64) *IndexedFieldInfo); ok {
		r0 = rf(fieldID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*IndexedFieldInfo)
		}
	}

	return r0
}

// MockSegment_GetIndex_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIndex'
type MockSegment_GetIndex_Call struct {
	*mock.Call
}

// GetIndex is a helper method to define mock.On call
//   - fieldID int64
func (_e *MockSegment_Expecter) GetIndex(fieldID interface{}) *MockSegment_GetIndex_Call {
	return &MockSegment_GetIndex_Call{Call: _e.mock.On("GetIndex", fieldID)}
}

func (_c *MockSegment_GetIndex_Call) Run(run func(fieldID int64)) *MockSegment_GetIndex_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockSegment_GetIndex_Call) Return(_a0 *IndexedFieldInfo) *MockSegment_GetIndex_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_GetIndex_Call) RunAndReturn(run func(int64) *IndexedFieldInfo) *MockSegment_GetIndex_Call {
	_c.Call.Return(run)
	return _c
}

// ID provides a mock function with given fields:
func (_m *MockSegment) ID() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockSegment_ID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ID'
type MockSegment_ID_Call struct {
	*mock.Call
}

// ID is a helper method to define mock.On call
func (_e *MockSegment_Expecter) ID() *MockSegment_ID_Call {
	return &MockSegment_ID_Call{Call: _e.mock.On("ID")}
}

func (_c *MockSegment_ID_Call) Run(run func()) *MockSegment_ID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_ID_Call) Return(_a0 int64) *MockSegment_ID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_ID_Call) RunAndReturn(run func() int64) *MockSegment_ID_Call {
	_c.Call.Return(run)
	return _c
}

// Indexes provides a mock function with given fields:
func (_m *MockSegment) Indexes() []*IndexedFieldInfo {
	ret := _m.Called()

	var r0 []*IndexedFieldInfo
	if rf, ok := ret.Get(0).(func() []*IndexedFieldInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*IndexedFieldInfo)
		}
	}

	return r0
}

// MockSegment_Indexes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Indexes'
type MockSegment_Indexes_Call struct {
	*mock.Call
}

// Indexes is a helper method to define mock.On call
func (_e *MockSegment_Expecter) Indexes() *MockSegment_Indexes_Call {
	return &MockSegment_Indexes_Call{Call: _e.mock.On("Indexes")}
}

func (_c *MockSegment_Indexes_Call) Run(run func()) *MockSegment_Indexes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_Indexes_Call) Return(_a0 []*IndexedFieldInfo) *MockSegment_Indexes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_Indexes_Call) RunAndReturn(run func() []*IndexedFieldInfo) *MockSegment_Indexes_Call {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function with given fields: rowIDs, timestamps, record
func (_m *MockSegment) Insert(rowIDs []int64, timestamps []uint64, record *segcorepb.InsertRecord) error {
	ret := _m.Called(rowIDs, timestamps, record)

	var r0 error
	if rf, ok := ret.Get(0).(func([]int64, []uint64, *segcorepb.InsertRecord) error); ok {
		r0 = rf(rowIDs, timestamps, record)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSegment_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type MockSegment_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - rowIDs []int64
//   - timestamps []uint64
//   - record *segcorepb.InsertRecord
func (_e *MockSegment_Expecter) Insert(rowIDs interface{}, timestamps interface{}, record interface{}) *MockSegment_Insert_Call {
	return &MockSegment_Insert_Call{Call: _e.mock.On("Insert", rowIDs, timestamps, record)}
}

func (_c *MockSegment_Insert_Call) Run(run func(rowIDs []int64, timestamps []uint64, record *segcorepb.InsertRecord)) *MockSegment_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]int64), args[1].([]uint64), args[2].(*segcorepb.InsertRecord))
	})
	return _c
}

func (_c *MockSegment_Insert_Call) Return(_a0 error) *MockSegment_Insert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_Insert_Call) RunAndReturn(run func([]int64, []uint64, *segcorepb.InsertRecord) error) *MockSegment_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// InsertCount provides a mock function with given fields:
func (_m *MockSegment) InsertCount() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockSegment_InsertCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertCount'
type MockSegment_InsertCount_Call struct {
	*mock.Call
}

// InsertCount is a helper method to define mock.On call
func (_e *MockSegment_Expecter) InsertCount() *MockSegment_InsertCount_Call {
	return &MockSegment_InsertCount_Call{Call: _e.mock.On("InsertCount")}
}

func (_c *MockSegment_InsertCount_Call) Run(run func()) *MockSegment_InsertCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_InsertCount_Call) Return(_a0 int64) *MockSegment_InsertCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_InsertCount_Call) RunAndReturn(run func() int64) *MockSegment_InsertCount_Call {
	_c.Call.Return(run)
	return _c
}

// LastDeltaTimestamp provides a mock function with given fields:
func (_m *MockSegment) LastDeltaTimestamp() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// MockSegment_LastDeltaTimestamp_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LastDeltaTimestamp'
type MockSegment_LastDeltaTimestamp_Call struct {
	*mock.Call
}

// LastDeltaTimestamp is a helper method to define mock.On call
func (_e *MockSegment_Expecter) LastDeltaTimestamp() *MockSegment_LastDeltaTimestamp_Call {
	return &MockSegment_LastDeltaTimestamp_Call{Call: _e.mock.On("LastDeltaTimestamp")}
}

func (_c *MockSegment_LastDeltaTimestamp_Call) Run(run func()) *MockSegment_LastDeltaTimestamp_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_LastDeltaTimestamp_Call) Return(_a0 uint64) *MockSegment_LastDeltaTimestamp_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_LastDeltaTimestamp_Call) RunAndReturn(run func() uint64) *MockSegment_LastDeltaTimestamp_Call {
	_c.Call.Return(run)
	return _c
}

// MayPkExist provides a mock function with given fields: pk
func (_m *MockSegment) MayPkExist(pk storage.PrimaryKey) bool {
	ret := _m.Called(pk)

	var r0 bool
	if rf, ok := ret.Get(0).(func(storage.PrimaryKey) bool); ok {
		r0 = rf(pk)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockSegment_MayPkExist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MayPkExist'
type MockSegment_MayPkExist_Call struct {
	*mock.Call
}

// MayPkExist is a helper method to define mock.On call
//   - pk storage.PrimaryKey
func (_e *MockSegment_Expecter) MayPkExist(pk interface{}) *MockSegment_MayPkExist_Call {
	return &MockSegment_MayPkExist_Call{Call: _e.mock.On("MayPkExist", pk)}
}

func (_c *MockSegment_MayPkExist_Call) Run(run func(pk storage.PrimaryKey)) *MockSegment_MayPkExist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(storage.PrimaryKey))
	})
	return _c
}

func (_c *MockSegment_MayPkExist_Call) Return(_a0 bool) *MockSegment_MayPkExist_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_MayPkExist_Call) RunAndReturn(run func(storage.PrimaryKey) bool) *MockSegment_MayPkExist_Call {
	_c.Call.Return(run)
	return _c
}

// MemSize provides a mock function with given fields:
func (_m *MockSegment) MemSize() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockSegment_MemSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MemSize'
type MockSegment_MemSize_Call struct {
	*mock.Call
}

// MemSize is a helper method to define mock.On call
func (_e *MockSegment_Expecter) MemSize() *MockSegment_MemSize_Call {
	return &MockSegment_MemSize_Call{Call: _e.mock.On("MemSize")}
}

func (_c *MockSegment_MemSize_Call) Run(run func()) *MockSegment_MemSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_MemSize_Call) Return(_a0 int64) *MockSegment_MemSize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_MemSize_Call) RunAndReturn(run func() int64) *MockSegment_MemSize_Call {
	_c.Call.Return(run)
	return _c
}

// Partition provides a mock function with given fields:
func (_m *MockSegment) Partition() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockSegment_Partition_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Partition'
type MockSegment_Partition_Call struct {
	*mock.Call
}

// Partition is a helper method to define mock.On call
func (_e *MockSegment_Expecter) Partition() *MockSegment_Partition_Call {
	return &MockSegment_Partition_Call{Call: _e.mock.On("Partition")}
}

func (_c *MockSegment_Partition_Call) Run(run func()) *MockSegment_Partition_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_Partition_Call) Return(_a0 int64) *MockSegment_Partition_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_Partition_Call) RunAndReturn(run func() int64) *MockSegment_Partition_Call {
	_c.Call.Return(run)
	return _c
}

// RLock provides a mock function with given fields:
func (_m *MockSegment) RLock() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSegment_RLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RLock'
type MockSegment_RLock_Call struct {
	*mock.Call
}

// RLock is a helper method to define mock.On call
func (_e *MockSegment_Expecter) RLock() *MockSegment_RLock_Call {
	return &MockSegment_RLock_Call{Call: _e.mock.On("RLock")}
}

func (_c *MockSegment_RLock_Call) Run(run func()) *MockSegment_RLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_RLock_Call) Return(_a0 error) *MockSegment_RLock_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_RLock_Call) RunAndReturn(run func() error) *MockSegment_RLock_Call {
	_c.Call.Return(run)
	return _c
}

// RUnlock provides a mock function with given fields:
func (_m *MockSegment) RUnlock() {
	_m.Called()
}

// MockSegment_RUnlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RUnlock'
type MockSegment_RUnlock_Call struct {
	*mock.Call
}

// RUnlock is a helper method to define mock.On call
func (_e *MockSegment_Expecter) RUnlock() *MockSegment_RUnlock_Call {
	return &MockSegment_RUnlock_Call{Call: _e.mock.On("RUnlock")}
}

func (_c *MockSegment_RUnlock_Call) Run(run func()) *MockSegment_RUnlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_RUnlock_Call) Return() *MockSegment_RUnlock_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockSegment_RUnlock_Call) RunAndReturn(run func()) *MockSegment_RUnlock_Call {
	_c.Call.Return(run)
	return _c
}

// RowNum provides a mock function with given fields:
func (_m *MockSegment) RowNum() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockSegment_RowNum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RowNum'
type MockSegment_RowNum_Call struct {
	*mock.Call
}

// RowNum is a helper method to define mock.On call
func (_e *MockSegment_Expecter) RowNum() *MockSegment_RowNum_Call {
	return &MockSegment_RowNum_Call{Call: _e.mock.On("RowNum")}
}

func (_c *MockSegment_RowNum_Call) Run(run func()) *MockSegment_RowNum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_RowNum_Call) Return(_a0 int64) *MockSegment_RowNum_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_RowNum_Call) RunAndReturn(run func() int64) *MockSegment_RowNum_Call {
	_c.Call.Return(run)
	return _c
}

// Shard provides a mock function with given fields:
func (_m *MockSegment) Shard() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockSegment_Shard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Shard'
type MockSegment_Shard_Call struct {
	*mock.Call
}

// Shard is a helper method to define mock.On call
func (_e *MockSegment_Expecter) Shard() *MockSegment_Shard_Call {
	return &MockSegment_Shard_Call{Call: _e.mock.On("Shard")}
}

func (_c *MockSegment_Shard_Call) Run(run func()) *MockSegment_Shard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_Shard_Call) Return(_a0 string) *MockSegment_Shard_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_Shard_Call) RunAndReturn(run func() string) *MockSegment_Shard_Call {
	_c.Call.Return(run)
	return _c
}

// StartPosition provides a mock function with given fields:
func (_m *MockSegment) StartPosition() *msgpb.MsgPosition {
	ret := _m.Called()

	var r0 *msgpb.MsgPosition
	if rf, ok := ret.Get(0).(func() *msgpb.MsgPosition); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*msgpb.MsgPosition)
		}
	}

	return r0
}

// MockSegment_StartPosition_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartPosition'
type MockSegment_StartPosition_Call struct {
	*mock.Call
}

// StartPosition is a helper method to define mock.On call
func (_e *MockSegment_Expecter) StartPosition() *MockSegment_StartPosition_Call {
	return &MockSegment_StartPosition_Call{Call: _e.mock.On("StartPosition")}
}

func (_c *MockSegment_StartPosition_Call) Run(run func()) *MockSegment_StartPosition_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_StartPosition_Call) Return(_a0 *msgpb.MsgPosition) *MockSegment_StartPosition_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_StartPosition_Call) RunAndReturn(run func() *msgpb.MsgPosition) *MockSegment_StartPosition_Call {
	_c.Call.Return(run)
	return _c
}

// Type provides a mock function with given fields:
func (_m *MockSegment) Type() commonpb.SegmentState {
	ret := _m.Called()

	var r0 commonpb.SegmentState
	if rf, ok := ret.Get(0).(func() commonpb.SegmentState); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(commonpb.SegmentState)
	}

	return r0
}

// MockSegment_Type_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Type'
type MockSegment_Type_Call struct {
	*mock.Call
}

// Type is a helper method to define mock.On call
func (_e *MockSegment_Expecter) Type() *MockSegment_Type_Call {
	return &MockSegment_Type_Call{Call: _e.mock.On("Type")}
}

func (_c *MockSegment_Type_Call) Run(run func()) *MockSegment_Type_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_Type_Call) Return(_a0 commonpb.SegmentState) *MockSegment_Type_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_Type_Call) RunAndReturn(run func() commonpb.SegmentState) *MockSegment_Type_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateBloomFilter provides a mock function with given fields: pks
func (_m *MockSegment) UpdateBloomFilter(pks []storage.PrimaryKey) {
	_m.Called(pks)
}

// MockSegment_UpdateBloomFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateBloomFilter'
type MockSegment_UpdateBloomFilter_Call struct {
	*mock.Call
}

// UpdateBloomFilter is a helper method to define mock.On call
//   - pks []storage.PrimaryKey
func (_e *MockSegment_Expecter) UpdateBloomFilter(pks interface{}) *MockSegment_UpdateBloomFilter_Call {
	return &MockSegment_UpdateBloomFilter_Call{Call: _e.mock.On("UpdateBloomFilter", pks)}
}

func (_c *MockSegment_UpdateBloomFilter_Call) Run(run func(pks []storage.PrimaryKey)) *MockSegment_UpdateBloomFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]storage.PrimaryKey))
	})
	return _c
}

func (_c *MockSegment_UpdateBloomFilter_Call) Return() *MockSegment_UpdateBloomFilter_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockSegment_UpdateBloomFilter_Call) RunAndReturn(run func([]storage.PrimaryKey)) *MockSegment_UpdateBloomFilter_Call {
	_c.Call.Return(run)
	return _c
}

// Version provides a mock function with given fields:
func (_m *MockSegment) Version() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MockSegment_Version_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Version'
type MockSegment_Version_Call struct {
	*mock.Call
}

// Version is a helper method to define mock.On call
func (_e *MockSegment_Expecter) Version() *MockSegment_Version_Call {
	return &MockSegment_Version_Call{Call: _e.mock.On("Version")}
}

func (_c *MockSegment_Version_Call) Run(run func()) *MockSegment_Version_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSegment_Version_Call) Return(_a0 int64) *MockSegment_Version_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSegment_Version_Call) RunAndReturn(run func() int64) *MockSegment_Version_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSegment creates a new instance of MockSegment. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSegment(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSegment {
	mock := &MockSegment{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
