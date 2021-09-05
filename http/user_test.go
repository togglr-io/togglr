package http_test

import (
	"bytes"
	"errors"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/togglr-io/togglr/http"
	"github.com/togglr-io/togglr/mock"
	"github.com/togglr-io/togglr/uid"
	"go.uber.org/zap"
)

func Test_HandleUserPost(t *testing.T) {
	accountID := uid.New().String()
	id := uid.New().String()
	cases := []struct {
		name                string
		payload             string
		userService         *mock.UserService
		expectedStatus      int
		expectedCreateCalls int
		expectedUpdateCalls int
	}{
		{
			name:                "successful create",
			payload:             fmt.Sprintf(`{"accountId": "%s", "email": "test@togglr.io"}`, accountID),
			userService:         mock.NewUserService(nil),
			expectedStatus:      200,
			expectedCreateCalls: 1,
			expectedUpdateCalls: 0,
		},
		{
			name:                "bad request",
			payload:             `{"accountId": invalid, "email": "test@togglr.io"}`,
			userService:         mock.NewUserService(nil),
			expectedStatus:      400,
			expectedCreateCalls: 0,
			expectedUpdateCalls: 0,
		},
		{
			name:                "failed create",
			payload:             fmt.Sprintf(`{"accountId": "%s", "email": "test@togglr.io"}`, accountID),
			userService:         mock.NewUserService(errors.New("forced")),
			expectedStatus:      500,
			expectedCreateCalls: 1,
			expectedUpdateCalls: 0,
		},
		{
			name:                "successful update",
			payload:             fmt.Sprintf(`{"id": "%s", "name": "Test User"}`, id),
			userService:         mock.NewUserService(nil),
			expectedStatus:      200,
			expectedCreateCalls: 0,
			expectedUpdateCalls: 1,
		},
		{
			name:                "failed update",
			payload:             fmt.Sprintf(`{"id": "%s", "name": "Test User"}`, id),
			userService:         mock.NewUserService(errors.New("forced")),
			expectedStatus:      500,
			expectedCreateCalls: 0,
			expectedUpdateCalls: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := http.Config{
				Logger: zap.NewNop(),
				Services: http.Services{
					UserService: c.userService,
				},
			}

			s := httptest.NewServer(http.BuildRoutes(cfg))
			defer s.Close()
			url := fmt.Sprintf("%s/user", s.URL)
			req, err := stdhttp.NewRequest("POST", url, bytes.NewReader([]byte(c.payload)))
			if err != nil {
				t.Fatalf("failed to create request: %s", err)
			}

			res, err := stdhttp.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %s", err)
			}

			if res.StatusCode != c.expectedStatus {
				t.Fatalf("expected status code of %d, but got %d", c.expectedStatus, res.StatusCode)
			}

			if c.userService.CreateUserCalled != c.expectedCreateCalls {
				t.Fatalf("expected CreateUser user to be called %d times, but it was called %d times", c.expectedCreateCalls, c.userService.CreateUserCalled)
			}

			if c.userService.UpdateUserCalled != c.expectedUpdateCalls {
				t.Fatalf("expected UpdateUser user to be called %d times, but it was called %d times", c.expectedUpdateCalls, c.userService.UpdateUserCalled)
			}
		})
	}
}

func Test_HandleUserGet(t *testing.T) {
	cases := []struct {
		name           string
		query          string
		userService    *mock.UserService
		expectedStatus int
		expectedCalls  int
	}{
		{
			name:           "successful test",
			query:          "",
			userService:    mock.NewUserService(nil),
			expectedStatus: 200,
			expectedCalls:  1,
		},
		{
			name:           "service failure",
			query:          "",
			userService:    mock.NewUserService(errors.New("forced")),
			expectedStatus: 500,
			expectedCalls:  1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := http.Config{
				Logger: zap.NewNop(),
				Services: http.Services{
					UserService: c.userService,
				},
			}

			s := httptest.NewServer(http.BuildRoutes(cfg))
			defer s.Close()
			url := fmt.Sprintf("%s/user?%s", s.URL, c.query)
			req, err := stdhttp.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("failed to create request: %s", err)
			}

			res, err := stdhttp.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %s", err)
			}

			if res.StatusCode != c.expectedStatus {
				t.Fatalf("expected status code of %d, but got %d", c.expectedStatus, res.StatusCode)
			}

			if c.userService.ListUsersCalled != c.expectedCalls {
				t.Fatalf("expected FetchUser user to be called %d times, but it was called %d times", c.expectedCalls, c.userService.CreateUserCalled)
			}
		})
	}
}

func Test_HandleUserDelete(t *testing.T) {
	cases := []struct {
		name           string
		id             string
		userService    *mock.UserService
		expectedStatus int
		expectedCalls  int
	}{
		{
			name:           "successful test",
			id:             "c149f08b-b0fa-4a5d-8a6c-03ac992aa454",
			userService:    mock.NewUserService(nil),
			expectedStatus: 204,
			expectedCalls:  1,
		},
		{
			name:           "bad request",
			id:             "123",
			userService:    mock.NewUserService(nil),
			expectedStatus: 400,
			expectedCalls:  0,
		},
		{
			name:           "service failure",
			id:             "c149f08b-b0fa-4a5d-8a6c-03ac992aa454",
			userService:    mock.NewUserService(errors.New("forced")),
			expectedStatus: 500,
			expectedCalls:  1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := http.Config{
				Logger: zap.NewNop(),
				Services: http.Services{
					UserService: c.userService,
				},
			}

			s := httptest.NewServer(http.BuildRoutes(cfg))
			defer s.Close()
			url := fmt.Sprintf("%s/user/%s", s.URL, c.id)
			req, err := stdhttp.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Fatalf("failed to create request: %s", err)
			}

			res, err := stdhttp.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %s", err)
			}

			if res.StatusCode != c.expectedStatus {
				t.Fatalf("expected status code of %d, but got %d", c.expectedStatus, res.StatusCode)
			}

			if c.userService.DeleteUserCalled != c.expectedCalls {
				t.Fatalf("expected DeleteUser to be called %d times, but it was called %d times", c.expectedCalls, c.userService.CreateUserCalled)
			}
		})
	}
}
