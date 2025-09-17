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
		expectedTimes  int
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
			expectedTimes: 1,
		},
		{
			name:           "error get person",
			inputID:        `{"id": 22}`,
			expectedPerson: nil,
			expectedError:  status.Error(codes.InvalidArgument, "error processing request"),
			expectedTimes:  1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGrpc := mock_httptests.NewMockHttpRequestGrpc(ctrl)

			var responsePerson struct {
				Id      int
				User_id int
				Name    string
				Admin   bool
			}
			json.Unmarshal([]byte(tc.inputID), &responsePerson)

			expectedReq := &g_serv.GetRequest{Id: int64(responsePerson.Id)}
			mockGrpc.EXPECT().
				GetRequestGrpc(gomock.Any(), expectedReq).
				Return(tc.expectedPerson, tc.expectedError).
				Times(tc.expectedTimes)

			server := &HttpServer{requestGrpc: mockGrpc}

			req := httptest.NewRequest("GET", "/", strings.NewReader(tc.inputID))
			w := httptest.NewRecorder()
			server.handleGet(w, req)

			if tc.expectedError != nil {
				assert.Empty(t, w.Body.String())
			} else {
				err := json.Unmarshal(w.Body.Bytes(), &responsePerson)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedPerson.Info.Name, responsePerson.Name)
				assert.Equal(t, int(tc.expectedPerson.Info.UserId), responsePerson.User_id)
				assert.Equal(t, tc.expectedPerson.Info.Admin, responsePerson.Admin)

			}
		})
	}
}
