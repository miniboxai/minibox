//go:generate mockgen -destination mock/clients.go minibox.ai/pkg/api/v1 ClientInterface

package v1

import (
	"context"

	types1 "github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
	"minibox.ai/pkg/api/v1/proto"
	"minibox.ai/pkg/api/v1/types"
)

type ClientInterface interface {
	proto.UserServiceClient
	proto.ProjectServiceClient
	proto.DatasetServiceClient
}

type Clients struct {
	token   string
	c       *grpc.ClientConn
	project proto.ProjectServiceClient
	user    proto.UserServiceClient
	dataset proto.DatasetServiceClient
}

func NewClients(c *grpc.ClientConn) *Clients {
	return &Clients{
		c:       c,
		project: proto.NewProjectServiceClient(c),
		user:    proto.NewUserServiceClient(c),
		dataset: proto.NewDatasetServiceClient(c),
	}
}

func (c *Clients) AddUser(ctx context.Context, in *types.User, opts ...grpc.CallOption) (*types1.Empty, error) {
	panic("not implemented")
}

func (c *Clients) ListUsers(ctx context.Context, in *types.ListUsersRequest, opts ...grpc.CallOption) (*types.UsersReply, error) {
	return c.user.ListUsers(ctx, in, opts...)
}

func (c *Clients) ListUsersByRole(ctx context.Context, in *types.UserRole, opts ...grpc.CallOption) (*types.UsersReply, error) {
	panic("not implemented")
}

func (c *Clients) UpdateUser(ctx context.Context, in *types.UpdateUserRequest, opts ...grpc.CallOption) (*types.User, error) {
	panic("not implemented")
}

func (c *Clients) ListProjects(ctx context.Context, in *types.ListProjectsRequest, opts ...grpc.CallOption) (*types.ProjectsReply, error) {
	return c.project.ListProjects(ctx, in, opts...)
}

func (c *Clients) CreateDataset(ctx context.Context, in *types.CreateDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error) {
	return c.dataset.CreateDataset(ctx, in, opts...)
}

func (c *Clients) ListDatasets(ctx context.Context, in *types.ListDatasetsRequest, opts ...grpc.CallOption) (*types.DatasetsReply, error) {
	return c.dataset.ListDatasets(ctx, in, opts...)
}

func (c *Clients) CreateProject(ctx context.Context, in *types.CreateProjectRequest, opts ...grpc.CallOption) (*types.Project, error) {
	return c.project.CreateProject(ctx, in, opts...)
}

func (c *Clients) GetProject(ctx context.Context, in *types.GetProjectRequest, opts ...grpc.CallOption) (*types.Project, error) {
	return c.project.GetProject(ctx, in, opts...)
}

func (c *Clients) DeleteProject(ctx context.Context, in *types.DeleteProjectRequest, opts ...grpc.CallOption) (*types.Project, error) {
	return c.project.DeleteProject(ctx, in, opts...)
}

func (c *Clients) UpdateProject(ctx context.Context, in *types.UpdateProjectRequest, opts ...grpc.CallOption) (*types.Project, error) {
	return c.project.UpdateProject(ctx, in, opts...)
}

func (c *Clients) DeleteDataset(ctx context.Context, in *types.DeleteDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error) {
	return c.dataset.DeleteDataset(ctx, in, opts...)
}

func (c *Clients) RestoreDataset(ctx context.Context, in *types.RestoreDatasetRequest, opts ...grpc.CallOption) (*types.Dataset, error) {
	return c.dataset.RestoreDataset(ctx, in, opts...)
}

func (c *Clients) ListDatasetObjects(ctx context.Context, in *types.ListObjectsRequest, opts ...grpc.CallOption) (*types.ObjectsReply, error) {
	return c.dataset.ListDatasetObjects(ctx, in, opts...)
}

func (c *Clients) GetDatasetObject(ctx context.Context, in *types.GetObjectRequest, opts ...grpc.CallOption) (proto.DatasetService_GetDatasetObjectClient, error) {
	return c.dataset.GetDatasetObject(ctx, in, opts...)
}

func (c *Clients) PutDatasetObject(ctx context.Context, opts ...grpc.CallOption) (proto.DatasetService_PutDatasetObjectClient, error) {
	return c.dataset.PutDatasetObject(ctx, opts...)
}

func (c *Clients) DeleteDatasetObject(ctx context.Context, in *types.DeleteObjectRequest, opts ...grpc.CallOption) (*types.ObjectReply, error) {
	panic("not implemented")
}
