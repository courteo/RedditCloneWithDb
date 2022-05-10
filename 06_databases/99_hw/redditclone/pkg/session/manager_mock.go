package session

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
)

type MockSessionRepo struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepoMockRecorder
}

type MockSessionRepoMockRecorder struct {
	mock *MockSessionRepo
}

func NewMockSessionRepo(ctrl *gomock.Controller) *MockSessionRepo {
	mock := &MockSessionRepo{ctrl: ctrl}
	mock.recorder = &MockSessionRepoMockRecorder{mock}
	return mock
}

func (m *MockSessionRepo) EXPECT() *MockSessionRepoMockRecorder {
	return m.recorder
}

func (m *MockSessionRepo) Check(login *http.Request) (*Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", login)
	ret0, _ := ret[0].(*Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockSessionRepoMockRecorder) Check(login *http.Request) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockSessionRepo)(nil).Check), login)
}

func (m *MockSessionRepo) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyCurrent", w, r)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockSessionRepoMockRecorder) DestroyCurrent(w http.ResponseWriter, r *http.Request) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyCurrent", reflect.TypeOf((*MockSessionRepo)(nil).DestroyCurrent), w, r)
}

func (m *MockSessionRepo) Create(w http.ResponseWriter, userID uint32, path string) (*Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", w, userID, path)
	ret0, _ := ret[0].(*Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockSessionRepoMockRecorder) Create(w http.ResponseWriter, userID uint32, path string) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionRepo)(nil).Create), w, userID, path)
}
