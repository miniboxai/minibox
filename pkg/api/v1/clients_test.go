package v1

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	// "github.com/golang/mock/sample"
	// "github.com/golang/mock/sample/imp1"
	mock "minibox.ai/pkg/api/v1/mock"
)

func TestClients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIndex := mock.NewMockClientInterface(ctrl)
	ctx := context.Background()
	mockIndex.ListUsers(ctx, nil) // literals work
	// mockIndex.EXPECT().ListUsers(ctx, nil) // literals work
	// mockIndex.EXPECT().ListUsers("b", gomock.Eq(2)) // matchers work too
	// mockIndex.EXPECT().NillableRet()
}
