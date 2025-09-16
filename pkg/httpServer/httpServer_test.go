package httpServer

import (
	"context"
	g_serv "grpc/pkg/proto"
	mock_httptests "grpc/testingHttp/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleGet(t *testing.T) {
	testCases := []struct {
		name           string
		inputID        int64
		expectedPerson *g_serv.GetResponse
		expectedError  error
		expectedTimes  int
	}{
		{
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
		},
		/*{
			name:           "error get person",
			inputID:        2,
			expectedPerson: nil,
			expectedError:  status.Error(codes.InvalidArgument, "error processing request"),
			expectedTimes:  0,
		},*/
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGrpc := mock_httptests.NewMockHttpRequestGrpc(ctrl)
			expectedReq := &g_serv.GetRequest{Id: tc.inputID}
			mockGrpc.EXPECT().
				GetRequestGrpc(gomock.Any(), expectedReq).
				Return(tc.expectedPerson, tc.expectedError).
				Times(tc.expectedTimes)

			server := &HttpServer{requestGrpc: mockGrpc}
			req := &g_serv.GetRequest{Id: tc.inputID}
			resp, err := server.requestGrpc.GetRequestGrpc(context.Background(), req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expectedPerson.Info.Name, resp.Info.Name)
			assert.Equal(t, tc.expectedPerson.Info.UserId, resp.Info.UserId)
			assert.Equal(t, tc.expectedPerson.Info.Admin, resp.Info.Admin)
		})
	}
}
