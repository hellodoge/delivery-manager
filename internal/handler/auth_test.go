package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/internal/service"
	mock_service "github.com/hellodoge/delivery-manager/internal/service/mocks"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user dm.User)

	const TestName = "John Doe"
	const TestUsername = "foobar"
	const TestPassword = "qwerty"
	const TestID = 1

	tests := []struct {
		name                 string
		inputBody            map[string]interface{}
		inputUser            dm.User
		behavior             mockBehavior
		expectedStatusCode   int
		expectedResponseBody *map[string]interface{}
	}{
		{
			name: "OK",
			inputBody: map[string]interface{}{
				"name":     TestName,
				"username": TestUsername,
				"password": TestPassword,
			},
			inputUser: dm.User{
				Name:     TestName,
				Username: TestUsername,
				Password: TestPassword,
			},
			behavior: func(s *mock_service.MockAuthorization, user dm.User) {
				s.EXPECT().CreateUser(user).Return(TestID, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: &map[string]interface{}{
				"id": TestID,
			},
		},
		{
			name: "Empty Name Field",
			inputBody: map[string]interface{}{
				"username": TestUsername,
				"password": TestPassword,
			},
			behavior:           func(s *mock_service.MockAuthorization, user dm.User) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Empty Username Field",
			inputBody: map[string]interface{}{
				"name":     TestName,
				"password": TestPassword,
			},
			behavior:           func(s *mock_service.MockAuthorization, user dm.User) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Empty Password Field",
			inputBody: map[string]interface{}{
				"name":     TestName,
				"username": TestUsername,
			},
			behavior:           func(s *mock_service.MockAuthorization, user dm.User) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Service Failure",
			inputBody: map[string]interface{}{
				"name":     TestName,
				"username": TestUsername,
				"password": TestPassword,
			},
			inputUser: dm.User{
				Name:     TestName,
				Username: TestUsername,
				Password: TestPassword,
			},
			behavior: func(s *mock_service.MockAuthorization, user dm.User) {
				s.EXPECT().CreateUser(user).Return(TestID, errors.New("internal error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// https://pkg.go.dev/github.com/golang/mock/gomock
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			test.behavior(auth, test.inputUser)

			services := &service.Service{Authorization: auth}
			handlers := Handler{services: services}

			const TestRoute = "/sign-up"

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.POST(TestRoute, handlers.signUp)

			requestBody, err := json.Marshal(test.inputBody)
			if err != nil {
				t.Error(err)
			}

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", TestRoute, bytes.NewBuffer(requestBody))
			router.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, test.expectedStatusCode)
			if test.expectedResponseBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				if err != nil {
					t.Error(err)
				}

				// https://tip.golang.org/doc/go1.12#fmt
				assert.Equal(t, fmt.Sprint(responseBody), fmt.Sprint(*test.expectedResponseBody))
			}
		})
	}
}
