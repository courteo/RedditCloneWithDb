package user

import (
	gomock "github.com/golang/mock/gomock"
	"reflect"
)

type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

func (m *MockUserRepo) NewUserID() uint32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewUserID")
	ret0, _ := ret[0].(uint32)
	return ret0
}

func (mr *MockUserRepoMockRecorder) NewUserID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewUserID", reflect.TypeOf((*MockUserRepo)(nil).NewUserID))
}

func (m *MockUserRepo) FindUser(login string) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", login)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUserRepoMockRecorder) FindUser(login string) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockUserRepo)(nil).FindUser), login)
}

func (m *MockUserRepo) Add(arg0 *User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockUserRepoMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUserRepo)(nil).Add), arg0)
}

func (m *MockUserRepo) Authorize(arg0 string, arg1 string) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", arg0, arg1)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUserRepoMockRecorder) Authorize(arg0 interface{}, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockUserRepo)(nil).Authorize), arg0, arg1)
}
