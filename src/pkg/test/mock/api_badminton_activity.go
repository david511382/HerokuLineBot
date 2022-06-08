// Code generated by MockGen. DO NOT EDIT.
// Source: ./src/logic/api/badminton_activity.go

// Package mock is a generated GoMock package.
package mock

import (
	error "heroku-line-bot/src/pkg/util/error"
	resp "heroku-line-bot/src/server/domain/resp"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockIBadmintonActivityApiLogic is a mock of IBadmintonActivityApiLogic interface.
type MockIBadmintonActivityApiLogic struct {
	ctrl     *gomock.Controller
	recorder *MockIBadmintonActivityApiLogicMockRecorder
}

// MockIBadmintonActivityApiLogicMockRecorder is the mock recorder for MockIBadmintonActivityApiLogic.
type MockIBadmintonActivityApiLogicMockRecorder struct {
	mock *MockIBadmintonActivityApiLogic
}

// NewMockIBadmintonActivityApiLogic creates a new mock instance.
func NewMockIBadmintonActivityApiLogic(ctrl *gomock.Controller) *MockIBadmintonActivityApiLogic {
	mock := &MockIBadmintonActivityApiLogic{ctrl: ctrl}
	mock.recorder = &MockIBadmintonActivityApiLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBadmintonActivityApiLogic) EXPECT() *MockIBadmintonActivityApiLogicMockRecorder {
	return m.recorder
}

// GetActivitys mocks base method.
func (m *MockIBadmintonActivityApiLogic) GetActivitys(fromDate, toDate *time.Time, pageIndex, pageSize uint, placeIDs, teamIDs []uint, everyWeekdays []time.Weekday) (resp.GetActivitys, error.IError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivitys", fromDate, toDate, pageIndex, pageSize, placeIDs, teamIDs, everyWeekdays)
	ret0, _ := ret[0].(resp.GetActivitys)
	ret1, _ := ret[1].(error.IError)
	return ret0, ret1
}

// GetActivitys indicates an expected call of GetActivitys.
func (mr *MockIBadmintonActivityApiLogicMockRecorder) GetActivitys(fromDate, toDate, pageIndex, pageSize, placeIDs, teamIDs, everyWeekdays interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivitys", reflect.TypeOf((*MockIBadmintonActivityApiLogic)(nil).GetActivitys), fromDate, toDate, pageIndex, pageSize, placeIDs, teamIDs, everyWeekdays)
}