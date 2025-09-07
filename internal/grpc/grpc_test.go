package my_grpc

import (
	"context"
	"grpc/pkg/database"
	g_serv "grpc/pkg/proto"
	mock_grpctests "grpc/testingGRPC/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGRPCServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_grpctests.NewMockDatabase(ctrl)

	expectedPerson := database.Person{
		Id:      1,
		Name:    "John",
		Admin:   true,
		User_id: 123,
	}

	mockDB.EXPECT().
		GetPerson(gomock.Any(), int64(1)).
		Return(&expectedPerson, nil).
		Times(1)

	server := &Server{Db: mockDB}
	req := &g_serv.GetRequest{Id: 1}
	resp, err := server.Get(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "John", resp.Info.Name)
	assert.Equal(t, int64(123), resp.Info.UserId)
	assert.True(t, resp.Info.Admin)

}
