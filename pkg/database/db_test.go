package database

/*func TestCreatePerson(t *testing.T) {
	testCases := []struct {
		name          string
		inputPerson   *g_serv.PostRequest
		expectedId    int64
		expectedError error
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
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_grpctests.NewMockDatabaseQS(ctrl)
			mockScan := mock_grpctests.NewMockDatabaseScan(ctrl)

			mockDB.EXPECT().
				QueryRow(gomock.Any(),
					"INSERT INTO people (User_id, Name, Admin) VALUES ($1, $2, $3) RETURNING id",
					tc.inputPerson.Info.UserId, tc.inputPerson.Info.Name, tc.inputPerson.Info.Admin).
				Return(mockScan)
			mockScan.EXPECT().
				Scan(gomock.Any()).
				Return(tc.expectedError)

			db := &Database{DB: mockDB}
		})
	}
}*/
