package httpServer

import (
	"encoding/json"
	g_serv "grpc/pkg/proto"
	mock_httptests "grpc/testingHttp/mock"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandleGet(t *testing.T) {
	testCases := []struct {
		name           string
		inputID        string
		expectedPerson *g_serv.GetResponse
		expectedError  error
	}{
		{
			name:    "successful get person",
			inputID: `{"id": 1}`,
			expectedPerson: &g_serv.GetResponse{
				Info: &g_serv.UserInfo{
					Name:   "John",
					Admin:  true,
					UserId: 123123,
				}},
			expectedError: nil,
		},
		{
			name:           "error get person",
			inputID:        `{"id": 22}`,
			expectedPerson: nil,
			expectedError:  status.Error(codes.InvalidArgument, "error processing request"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGrpc := mock_httptests.NewMockHttpRequestGrpc(ctrl)

			var person struct {
				Id      int64
				User_id int64
				Name    string
				Admin   bool
			}
			json.Unmarshal([]byte(tc.inputID), &person)

			mockReq := &g_serv.GetRequest{Id: person.Id}
			mockGrpc.EXPECT().
				GetRequestGrpc(gomock.Any(), mockReq).
				Return(tc.expectedPerson, tc.expectedError)

			server := &HttpServer{requestGrpc: mockGrpc}

			req := httptest.NewRequest("GET", "/", strings.NewReader(tc.inputID))
			w := httptest.NewRecorder()
			server.handleGet(w, req)

			if tc.expectedError != nil {
				assert.Empty(t, w.Body.String())
			} else {
				err := json.Unmarshal(w.Body.Bytes(), &person)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedPerson.Info.Name, person.Name)
				assert.Equal(t, tc.expectedPerson.Info.UserId, person.User_id)
				assert.Equal(t, tc.expectedPerson.Info.Admin, person.Admin)

			}
		})
	}
}

func TestHandlePost(t *testing.T) {
	testCases := []struct {
		name          string
		inputPerson   string
		expectedId    *g_serv.PostResponse
		expectedError error
	}{
		{
			name:          "successful post person",
			inputPerson:   `{"User_id": 986151, "Name": "Dima", "Admin": true}`,
			expectedId:    &g_serv.PostResponse{Id: 8},
			expectedError: nil,
		},
		{
			name:          "error post person",
			inputPerson:   `{"User_id": 127, "Name": "Anna", "Admin": false}`,
			expectedId:    nil,
			expectedError: status.Error(codes.InvalidArgument, "error processing request"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGrpc := mock_httptests.NewMockHttpRequestGrpc(ctrl)

			var person struct {
				Id      int64
				User_id int64
				Name    string
				Admin   bool
			}
			err := json.Unmarshal([]byte(tc.inputPerson), &person)
			assert.NoError(t, err)

			mockReq := &g_serv.PostRequest{
				Info: &g_serv.UserInfo{
					UserId: person.User_id,
					Name:   person.Name,
					Admin:  person.Admin,
				},
			}
			mockGrpc.EXPECT().
				PostRequestGrpc(gomock.Any(), mockReq).
				Return(tc.expectedId, tc.expectedError)
			server := &HttpServer{requestGrpc: mockGrpc}

			writer := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/", strings.NewReader(tc.inputPerson))
			server.handlePost(writer, request)

			if tc.expectedError != nil {
				assert.Empty(t, writer.Body.String())
			} else {
				err := json.Unmarshal(writer.Body.Bytes(), &person)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedId.Id, person.Id)

			}
		})
	}
}
