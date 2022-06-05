// Code generated by MockGen. DO NOT EDIT.
// Source: ./src/logic/badminton/activity_logic.go

// Package mock is a generated GoMock package.
package mock

import (
	badminton "heroku-line-bot/src/logic/badminton"
	error "heroku-line-bot/src/pkg/util/error"
	flow "heroku-line-bot/src/pkg/util/flow"
	activity "heroku-line-bot/src/repo/database/database/clubdb/activity"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockIBadmintonActivityLogic is a mock of IBadmintonActivityLogic interface.
type MockIBadmintonActivityLogic struct {
	ctrl     *gomock.Controller
	recorder *MockIBadmintonActivityLogicMockRecorder
}

// MockIBadmintonActivityLogicMockRecorder is the mock recorder for MockIBadmintonActivityLogic.
type MockIBadmintonActivityLogicMockRecorder struct {
	mock *MockIBadmintonActivityLogic
}

// NewMockIBadmintonActivityLogic creates a new mock instance.
func NewMockIBadmintonActivityLogic(ctrl *gomock.Controller) *MockIBadmintonActivityLogic {
	mock := &MockIBadmintonActivityLogic{ctrl: ctrl}
	mock.recorder = &MockIBadmintonActivityLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBadmintonActivityLogic) EXPECT() *MockIBadmintonActivityLogicMockRecorder {
	return m.recorder
}

// GetActivityDetail mocks base method.
func (m *MockIBadmintonActivityLogic) GetActivityDetail(activityReqs *activity.Reqs, respActivityID_detailsMap map[uint][]*badminton.CourtDetail) flow.IStep {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivityDetail", activityReqs, respActivityID_detailsMap)
	ret0, _ := ret[0].(flow.IStep)
	return ret0
}

// GetActivityDetail indicates an expected call of GetActivityDetail.
func (mr *MockIBadmintonActivityLogicMockRecorder) GetActivityDetail(activityReqs, respActivityID_detailsMap interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivityDetail", reflect.TypeOf((*MockIBadmintonActivityLogic)(nil).GetActivityDetail), activityReqs, respActivityID_detailsMap)
}

// GetUnfinishedActiviysSqlReqs mocks base method.
func (m *MockIBadmintonActivityLogic) GetUnfinishedActiviysSqlReqs(fromDate, toDate *time.Time, teamIDs, placeIDs []uint, everyWeekdays []time.Weekday) ([]*activity.Reqs, error.IError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnfinishedActiviysSqlReqs", fromDate, toDate, teamIDs, placeIDs, everyWeekdays)
	ret0, _ := ret[0].([]*activity.Reqs)
	ret1, _ := ret[1].(error.IError)
	return ret0, ret1
}

// GetUnfinishedActiviysSqlReqs indicates an expected call of GetUnfinishedActiviysSqlReqs.
func (mr *MockIBadmintonActivityLogicMockRecorder) GetUnfinishedActiviysSqlReqs(fromDate, toDate, teamIDs, placeIDs, everyWeekdays interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnfinishedActiviysSqlReqs", reflect.TypeOf((*MockIBadmintonActivityLogic)(nil).GetUnfinishedActiviysSqlReqs), fromDate, toDate, teamIDs, placeIDs, everyWeekdays)
}
