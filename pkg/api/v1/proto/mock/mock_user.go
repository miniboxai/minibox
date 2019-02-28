// Code generated by MockGen. DO NOT EDIT.
// Source: minibox.ai/pkg/api/v1/proto (interfaces: UserServiceClient)

// Package mock_proto is a generated GoMock package.
package mock_proto

import (
	context "context"
	types "github.com/gogo/protobuf/types"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	types0 "minibox.ai/pkg/api/v1/types"
	reflect "reflect"
)

// MockUserServiceClient is a mock of UserServiceClient interface
type MockUserServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceClientMockRecorder
}

// MockUserServiceClientMockRecorder is the mock recorder for MockUserServiceClient
type MockUserServiceClientMockRecorder struct {
	mock *MockUserServiceClient
}

// NewMockUserServiceClient creates a new mock instance
func NewMockUserServiceClient(ctrl *gomock.Controller) *MockUserServiceClient {
	mock := &MockUserServiceClient{ctrl: ctrl}
	mock.recorder = &MockUserServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserServiceClient) EXPECT() *MockUserServiceClientMockRecorder {
	return m.recorder
}

// AddUser mocks base method
func (m *MockUserServiceClient) AddUser(arg0 context.Context, arg1 *types0.User, arg2 ...grpc.CallOption) (*types.Empty, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddUser", varargs...)
	ret0, _ := ret[0].(*types.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser
func (mr *MockUserServiceClientMockRecorder) AddUser(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockUserServiceClient)(nil).AddUser), varargs...)
}

// ListUsers mocks base method
func (m *MockUserServiceClient) ListUsers(arg0 context.Context, arg1 *types0.ListUsersRequest, arg2 ...grpc.CallOption) (*types0.UsersReply, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListUsers", varargs...)
	ret0, _ := ret[0].(*types0.UsersReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers
func (mr *MockUserServiceClientMockRecorder) ListUsers(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockUserServiceClient)(nil).ListUsers), varargs...)
}

// ListUsersByRole mocks base method
func (m *MockUserServiceClient) ListUsersByRole(arg0 context.Context, arg1 *types0.UserRole, arg2 ...grpc.CallOption) (*types0.UsersReply, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListUsersByRole", varargs...)
	ret0, _ := ret[0].(*types0.UsersReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsersByRole indicates an expected call of ListUsersByRole
func (mr *MockUserServiceClientMockRecorder) ListUsersByRole(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsersByRole", reflect.TypeOf((*MockUserServiceClient)(nil).ListUsersByRole), varargs...)
}

// UpdateUser mocks base method
func (m *MockUserServiceClient) UpdateUser(arg0 context.Context, arg1 *types0.UpdateUserRequest, arg2 ...grpc.CallOption) (*types0.User, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateUser", varargs...)
	ret0, _ := ret[0].(*types0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser
func (mr *MockUserServiceClientMockRecorder) UpdateUser(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserServiceClient)(nil).UpdateUser), varargs...)
}
