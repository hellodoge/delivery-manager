package service

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/hellodoge/delivery-manager/dm"
	mock_repository "github.com/hellodoge/delivery-manager/internal/repository/mocks"
	"github.com/hellodoge/delivery-manager/pkg/response"
	"github.com/magiconair/properties/assert"
	"testing"
)

type userMatcher struct {
	user dm.User
}

func newUserMatcher(user dm.User) gomock.Matcher {
	return &userMatcher{user: user}
}

func (m *userMatcher) Matches(x interface{}) bool {
	user, ok := x.(dm.User)
	if !ok {
		return false
	}
	if user.Name != m.user.Name || user.Username != m.user.Username || user.Password != m.user.Password {
		return false
	}
	if user.PasswordHash == "" || user.PasswordSalt == "" {
		return false
	}
	return true
}

func (m *userMatcher) String() string {
	return fmt.Sprint(m.user)
}

func TestAuthService_CreateUser(t *testing.T) {
	type mockBehaviour func(r *mock_repository.MockAuthorization, user dm.User)

	const TestName = "John Doe"
	const TestUsername = "foobar_1"
	const TestPassword = "qwerty"
	const TestID = 1

	tests := []struct {
		name              string
		user              dm.User
		wantErrorParams   bool
		wantInternalError bool
		mockBehaviour     mockBehaviour
		id                int
	}{
		{
			name: "OK",
			user: dm.User{
				Name:     TestName,
				Username: TestUsername,
				Password: TestPassword,
			},
			mockBehaviour: func(r *mock_repository.MockAuthorization, user dm.User) {
				r.EXPECT().CreateUser(newUserMatcher(user)).Return(TestID, nil)
			},
			id: TestID,
		},
		{
			name: "Invalid username (space inside)",
			user: dm.User{
				Name:     TestName,
				Username: "foo bar",
				Password: TestPassword,
			},
			mockBehaviour:   func(r *mock_repository.MockAuthorization, user dm.User) {},
			wantErrorParams: true,
		},
		{
			name: "Invalid username (invalid character)",
			user: dm.User{
				Name:     TestName,
				Username: "foo$bar",
				Password: TestPassword,
			},
			mockBehaviour:   func(r *mock_repository.MockAuthorization, user dm.User) {},
			wantErrorParams: true,
		},
		{
			name: "Invalid username (starts with digit)",
			user: dm.User{
				Name:     TestName,
				Username: "1foobar",
				Password: TestPassword,
			},
			mockBehaviour:   func(r *mock_repository.MockAuthorization, user dm.User) {},
			wantErrorParams: true,
		},
		{
			name: "Invalid username (too short)",
			user: dm.User{
				Name:     TestName,
				Username: "foo",
				Password: TestPassword,
			},
			mockBehaviour:   func(r *mock_repository.MockAuthorization, user dm.User) {},
			wantErrorParams: true,
		},
		{
			name: "Too short password",
			user: dm.User{
				Name:     TestName,
				Username: TestUsername,
				Password: "qwe",
			},
			mockBehaviour:   func(r *mock_repository.MockAuthorization, user dm.User) {},
			wantErrorParams: true,
		},
		{
			name: "Repository Error",
			user: dm.User{
				Name:     TestName,
				Username: TestUsername,
				Password: TestPassword,
			},
			mockBehaviour: func(r *mock_repository.MockAuthorization, user dm.User) {
				r.EXPECT().CreateUser(newUserMatcher(user)).Return(0, errors.New("internal error"))
			},
			wantInternalError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authRepo := mock_repository.NewMockAuthorization(c)
			test.mockBehaviour(authRepo, test.user)

			auth := AuthService{repo: authRepo}
			id, err := auth.CreateUser(test.user)
			if err != nil {
				if test.wantErrorParams {
					if _, ok := err.(response.ErrorResponseParameters); !ok {
						t.Error("wanted ErrorResponseParameters, got Internal error")
					}
				} else if !test.wantInternalError {
					t.Error(err)
				}
			} else if test.wantInternalError || test.wantErrorParams {
				t.Error("wanted error")
			} else {
				assert.Equal(t, id, test.id)
			}
		})
	}
}
