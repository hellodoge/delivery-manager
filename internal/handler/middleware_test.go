package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/hellodoge/delivery-manager/internal/service"
	mock_service "github.com/hellodoge/delivery-manager/internal/service/mocks"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockAuthorization, token string)

	const TestToken = "foobar"
	const TestID = 1

	const ResponseIDKey = "id"

	tests := []struct {
		name                 string
		headers              map[string]string
		token                string
		behaviour            mockBehaviour
		expectedStatusCode   int
		expectedResponseBody *map[string]interface{}
	}{
		{
			name: "OK",
			headers: map[string]string{
				"Authorization": "Bearer " + TestToken,
			},
			token: TestToken,
			behaviour: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(TestToken).Return(TestID, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: &map[string]interface{}{
				ResponseIDKey: TestID,
			},
		},
		{
			name:               "No Authorization header",
			behaviour:          func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Bearer misspelled",
			headers: map[string]string{
				"Authorization": "Bear " + TestToken,
			},
			token:              TestToken,
			behaviour:          func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Empty Bearer token",
			headers: map[string]string{
				"Authorization": "Bearer " + "",
			},
			token:              TestToken,
			behaviour:          func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid Bearer token",
			headers: map[string]string{
				"Authorization": "Bearer " + TestToken,
			},
			token: TestToken,
			behaviour: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(0, response.ErrorResponseParameters{
					Message:    "Invalid token",
					StatusCode: http.StatusUnauthorized,
				})
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Service error",
			headers: map[string]string{
				"Authorization": "Bearer " + TestToken,
			},
			token: TestToken,
			behaviour: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(0, response.ErrorResponseParameters{
					IsInternal: true,
				})
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
			test.behaviour(auth, test.token)

			services := &service.Service{Authorization: auth}
			handlers := Handler{services: services}

			const TestRoute = "/user-identity"

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			router.POST(TestRoute, handlers.userIdentity, func(c *gin.Context) {
				if id, ok := c.Get(userIdContextKey); ok {
					c.JSON(http.StatusOK, map[string]interface{}{
						ResponseIDKey: id,
					})
				} else {
					t.Error("User id not found in context")
				}
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", TestRoute, nil)
			for key, value := range test.headers {
				request.Header.Set(key, value)
			}

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
