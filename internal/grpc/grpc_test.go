package my_grpc

import (
	"context"
	"grpc/pkg/database"
	g_serv "grpc/pkg/proto"
	mock_grpctests "grpc/testingGRPC/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCServer_Get(t *testing.T) {
	testCases := []struct {
		name           string
		inputID        int64
		expectedPerson *database.Person
		expectedError  error
		expectedTimes  int
	}{
		{
			name:    "successful get person",
			inputID: 1,
			expectedPerson: &database.Person{
				Id:      1,
				Name:    "John",
				Admin:   true,
				User_id: 123123,
			},
			expectedError: nil,
			expectedTimes: 1,
		},
		{
			name:           "person not found",
			inputID:        999,
			expectedPerson: nil,
			expectedError:  status.Error(codes.InvalidArgument, "empty field"),
			expectedTimes:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDB := mock_grpctests.NewMockDatabase(ctrl)

			mockDB.EXPECT().
				GetPerson(gomock.Any(), tc.inputID).
				Return(tc.expectedPerson, tc.expectedError).
				Times(tc.expectedTimes)

			server := &Server{Db: mockDB}
			req := &g_serv.GetRequest{Id: tc.inputID}
			resp, err := server.Get(context.Background(), req)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.expectedPerson.Name, resp.Info.Name)
				assert.EqualValues(t, tc.expectedPerson.User_id, resp.Info.UserId)
				assert.Equal(t, tc.expectedPerson.Admin, resp.Info.Admin)
			}
		})
	}
}
func TestGRPCServer_Post(t *testing.T) {
	testCases := []struct {
		name          string
		inputPerson   *g_serv.PostRequest
		expectedId    int64
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
			expectedId:    12,
			expectedError: nil,
			expectedTimes: 1,
		},
		{
			name: "empty field",
			inputPerson: &g_serv.PostRequest{
				Info: &g_serv.UserInfo{
					Name:   "",
					Admin:  false,
					UserId: 154354,
				}},
			expectedId:    0,
			expectedError: status.Error(codes.InvalidArgument, "empty field"),
			expectedTimes: 0,
		},
		{
			name: "incorrect id lenght",
			inputPerson: &g_serv.PostRequest{
				Info: &g_serv.UserInfo{
					Name:   "Dima",
					Admin:  false,
					UserId: 13,
				}},
			expectedId:    0,
			expectedError: status.Error(codes.InvalidArgument, "incorrect user_id lenght"),
			expectedTimes: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_grpctests.NewMockDatabase(ctrl)
			mockDB.EXPECT().
				CreatePerson(gomock.Any(), tc.inputPerson).
				Return(tc.expectedId, tc.expectedError).
				Times(tc.expectedTimes)

			server := &Server{Db: mockDB}
			req := tc.inputPerson
			resp, err := server.Post(context.Background(), req)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.expectedId, resp.Id)
			}
		})
	}
}
