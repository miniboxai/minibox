// Code generated by MockGen. DO NOT EDIT.
// Source: minibox.ai/pkg/api/v1 (interfaces: ClientInterface)

// Package mock_v1 is a generated GoMock package.
package mock_v1

import (
	context "context"
	types "github.com/gogo/protobuf/types"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	types0 "minibox.ai/pkg/api/v1/types"
	reflect "reflect"
)

// MockClientInterface is a mock of ClientInterface interface
type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
}

// MockClientInterfaceMockRecorder is the mock recorder for MockClientInterface
type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

// NewMockClientInterface creates a new mock instance
func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

// AddUser mocks base method
func (m *MockClientInterface) AddUser(arg0 context.Context, arg1 *types0.User, arg2 ...grpc.CallOption) (*types.Empty, error) {
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
func (mr *MockClientInterfaceMockRecorder) AddUser(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockClientInterface)(nil).AddUser), varargs...)
}

// CreateDataset mocks base method
func (m *MockClientInterface) CreateDataset(arg0 context.Context, arg1 *types0.CreateDatasetRequest, arg2 ...grpc.CallOption) (*types0.Dataset, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateDataset", varargs...)
	ret0, _ := ret[0].(*types0.Dataset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDataset indicates an expected call of CreateDataset
func (mr *MockClientInterfaceMockRecorder) CreateDataset(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDataset", reflect.TypeOf((*MockClientInterface)(nil).CreateDataset), varargs...)
}

// CreateProject mocks base method
func (m *MockClientInterface) CreateProject(arg0 context.Context, arg1 *types0.CreateProjectRequest, arg2 ...grpc.CallOption) (*types0.Project, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateProject", varargs...)
	ret0, _ := ret[0].(*types0.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject
func (mr *MockClientInterfaceMockRecorder) CreateProject(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockClientInterface)(nil).CreateProject), varargs...)
}

// DeleteDataset mocks base method
func (m *MockClientInterface) DeleteDataset(arg0 context.Context, arg1 *types0.DeleteDatasetRequest, arg2 ...grpc.CallOption) (*types0.Dataset, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteDataset", varargs...)
	ret0, _ := ret[0].(*types0.Dataset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteDataset indicates an expected call of DeleteDataset
func (mr *MockClientInterfaceMockRecorder) DeleteDataset(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDataset", reflect.TypeOf((*MockClientInterface)(nil).DeleteDataset), varargs...)
}

// DeleteProject mocks base method
func (m *MockClientInterface) DeleteProject(arg0 context.Context, arg1 *types0.DeleteProjectRequest, arg2 ...grpc.CallOption) (*types0.Project, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteProject", varargs...)
	ret0, _ := ret[0].(*types0.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteProject indicates an expected call of DeleteProject
func (mr *MockClientInterfaceMockRecorder) DeleteProject(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockClientInterface)(nil).DeleteProject), varargs...)
}

// GetProject mocks base method
func (m *MockClientInterface) GetProject(arg0 context.Context, arg1 *types0.GetProjectRequest, arg2 ...grpc.CallOption) (*types0.Project, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProject", varargs...)
	ret0, _ := ret[0].(*types0.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject
func (mr *MockClientInterfaceMockRecorder) GetProject(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockClientInterface)(nil).GetProject), varargs...)
}

// ListDatasets mocks base method
func (m *MockClientInterface) ListDatasets(arg0 context.Context, arg1 *types0.ListDatasetsRequest, arg2 ...grpc.CallOption) (*types0.DatasetsReply, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListDatasets", varargs...)
	ret0, _ := ret[0].(*types0.DatasetsReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDatasets indicates an expected call of ListDatasets
func (mr *MockClientInterfaceMockRecorder) ListDatasets(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDatasets", reflect.TypeOf((*MockClientInterface)(nil).ListDatasets), varargs...)
}

// ListObjects mocks base method
func (m *MockClientInterface) ListObjects(arg0 context.Context, arg1 *types0.ListObjectsRequest, arg2 ...grpc.CallOption) (*types0.ObjectsReply, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListObjects", varargs...)
	ret0, _ := ret[0].(*types0.ObjectsReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListObjects indicates an expected call of ListObjects
func (mr *MockClientInterfaceMockRecorder) ListObjects(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListObjects", reflect.TypeOf((*MockClientInterface)(nil).ListObjects), varargs...)
}

// ListProject mocks base method
func (m *MockClientInterface) ListProject(arg0 context.Context, arg1 *types0.User, arg2 ...grpc.CallOption) (*types0.ProjectsReply, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListProject", varargs...)
	ret0, _ := ret[0].(*types0.ProjectsReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProject indicates an expected call of ListProject
func (mr *MockClientInterfaceMockRecorder) ListProject(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProject", reflect.TypeOf((*MockClientInterface)(nil).ListProject), varargs...)
}

// ListUsers mocks base method
func (m *MockClientInterface) ListUsers(arg0 context.Context, arg1 *types0.ListUsersRequest, arg2 ...grpc.CallOption) (*types0.UsersReply, error) {
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
func (mr *MockClientInterfaceMockRecorder) ListUsers(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockClientInterface)(nil).ListUsers), varargs...)
}

// ListUsersByRole mocks base method
func (m *MockClientInterface) ListUsersByRole(arg0 context.Context, arg1 *types0.UserRole, arg2 ...grpc.CallOption) (*types0.UsersReply, error) {
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
func (mr *MockClientInterfaceMockRecorder) ListUsersByRole(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsersByRole", reflect.TypeOf((*MockClientInterface)(nil).ListUsersByRole), varargs...)
}

// RestoreDataset mocks base method
func (m *MockClientInterface) RestoreDataset(arg0 context.Context, arg1 *types0.RestoreDatasetRequest, arg2 ...grpc.CallOption) (*types0.Dataset, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RestoreDataset", varargs...)
	ret0, _ := ret[0].(*types0.Dataset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestoreDataset indicates an expected call of RestoreDataset
func (mr *MockClientInterfaceMockRecorder) RestoreDataset(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreDataset", reflect.TypeOf((*MockClientInterface)(nil).RestoreDataset), varargs...)
}

// UpdateProject mocks base method
func (m *MockClientInterface) UpdateProject(arg0 context.Context, arg1 *types0.UpdateProjectRequest, arg2 ...grpc.CallOption) (*types0.Project, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateProject", varargs...)
	ret0, _ := ret[0].(*types0.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProject indicates an expected call of UpdateProject
func (mr *MockClientInterfaceMockRecorder) UpdateProject(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProject", reflect.TypeOf((*MockClientInterface)(nil).UpdateProject), varargs...)
}

// UpdateUser mocks base method
func (m *MockClientInterface) UpdateUser(arg0 context.Context, arg1 *types0.UpdateUserRequest, arg2 ...grpc.CallOption) (*types0.User, error) {
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
func (mr *MockClientInterfaceMockRecorder) UpdateUser(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockClientInterface)(nil).UpdateUser), varargs...)
}