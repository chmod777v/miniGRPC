package grpcconect

import (
	"context"
	g_serv "grpc/pkg/proto"
	mock_httptests "grpc/testingHttp/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRequestGrpc_Get(t *testing.T) {
	testCases := []struct {
		name           string
		inputID        int64
		expectedPerson *g_serv.GetResponse
		expectedError  error
		expectedTimes  int
	}{{
		name:    "successful get person",
		inputID: 1,
		expectedPerson: &g_serv.GetResponse{
			Info: &g_serv.UserInfo{
				Name:   "John",
				Admin:  true,
				UserId: 123123,
			}},
		expectedError: nil,
		expectedTimes: 1,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockGrpc := mock_httptests.NewMockHTTPServer(ctrl)

			expectedReq := &g_serv.GetRequest{Id: tc.inputID}
			mockGrpc.EXPECT().
				Get(gomock.Any(), expectedReq).
				Return(tc.expectedPerson, tc.expectedError).
				Times(tc.expectedTimes)

			server := &Server{grpcClient: mockGrpc}
			req := &g_serv.GetRequest{Id: tc.inputID}
			resp, err := server.grpcClient.Get(context.Background(), req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expectedPerson.Info.Name, resp.Info.Name)
			assert.Equal(t, tc.expectedPerson.Info.UserId, resp.Info.UserId)
			assert.Equal(t, tc.expectedPerson.Info.Admin, resp.Info.Admin)

		})
	}
}

func TestRequestGrpc_Post(t *testing.T) {
	testCases := []struct {
		name          string
		inputPerson   *g_serv.PostRequest
		expectedId    *g_serv.PostResponse
		expectedError error
		expectedTimes int
	}{
		{
			name: "successful post person",
			inputPerson: &g_serv.PostRequest{
				Info: &g_serv.UserInfo{
					Name:   "Ivan",
					Admin:  true,
					UserId: 123123,
				}},
			expectedId:    &g_serv.PostResponse{Id: 12},
			expectedError: nil,
			expectedTimes: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockGrpc := mock_httptests.NewMockHTTPServer(ctrl)

			mockGrpc.EXPECT().
				Post(gomock.Any(), tc.inputPerson).
				Return(tc.expectedId, tc.expectedError).
				Times(tc.expectedTimes)

			server := &Server{grpcClient: mockGrpc}
			req := tc.inputPerson
			resp, err := server.grpcClient.Post(context.Background(), req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expectedId.Id, resp.Id)

		})
	}
}
