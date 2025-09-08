package my_grpc

import (
	"context"
	"errors"
	"grpc/pkg/database"
	g_serv "grpc/pkg/proto"
	mock_grpctests "grpc/testingGRPC/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
				User_id: 123,
			},
			expectedError: nil,
			expectedTimes: 1,
		},
		{
			name:           "person not found",
			inputID:        999,
			expectedPerson: nil,
			expectedError:  errors.New("person not found"),
			expectedTimes:  1,
		},
		{
			name:    "another successful case",
			inputID: 2,
			expectedPerson: &database.Person{
				Id:      2,
				Name:    "Alice",
				Admin:   false,
				User_id: 456,
			},
			expectedError: nil,
			expectedTimes: 1,
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
