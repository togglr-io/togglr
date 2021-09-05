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

func Test_HandleAccountPOST(t *testing.T) {
	id := uid.New().String()
	cases := []struct {
		name                string
		payload             string
		accountService      *mock.AccountService
		expectedStatus      int
		expectedCreateCalls int
		expectedUpdateCalls int
	}{
		{
			name:                "successful create",
			payload:             `{"name": "Test Account"}`,
			accountService:      mock.NewAccountService(nil),
			expectedStatus:      200,
			expectedCreateCalls: 1,
			expectedUpdateCalls: 0,
		},
		{
			name:                "bad request",
			payload:             `{"name": invalid"Test Account"}`,
			accountService:      mock.NewAccountService(nil),
			expectedStatus:      400,
			expectedCreateCalls: 0,
			expectedUpdateCalls: 0,
		},
		{
			name:                "failed create",
			payload:             `{"name": "Test Account"}`,
			accountService:      mock.NewAccountService(errors.New("forced")),
			expectedStatus:      500,
			expectedCreateCalls: 1,
			expectedUpdateCalls: 0,
		},
		{
			name:                "successful update",
			payload:             fmt.Sprintf(`{"id": "%s", "name": "New Account"}`, id),
			accountService:      mock.NewAccountService(nil),
			expectedStatus:      200,
			expectedCreateCalls: 0,
			expectedUpdateCalls: 1,
		},
		{
			name:                "failed update",
			payload:             fmt.Sprintf(`{"id": "%s", "name": "New Account"}`, id),
			accountService:      mock.NewAccountService(errors.New("forced")),
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
					AccountService: c.accountService,
				},
			}

			s := httptest.NewServer(http.BuildRoutes(cfg))
			defer s.Close()
			url := fmt.Sprintf("%s/account", s.URL)
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

			if c.accountService.CreateAccountCalled != c.expectedCreateCalls {
				t.Fatalf("expected CreateAccount account to be called %d times, but it was called %d times", c.expectedCreateCalls, c.accountService.CreateAccountCalled)
			}

			if c.accountService.UpdateAccountCalled != c.expectedUpdateCalls {
				t.Fatalf("expected UpdateAccount account to be called %d times, but it was called %d times", c.expectedUpdateCalls, c.accountService.UpdateAccountCalled)
			}
		})
	}
}

func Test_HandleAccountGet(t *testing.T) {
	cases := []struct {
		name           string
		query          string
		accountService *mock.AccountService
		expectedStatus int
		expectedCalls  int
	}{
		{
			name:           "successful test",
			query:          "",
			accountService: mock.NewAccountService(nil),
			expectedStatus: 200,
			expectedCalls:  1,
		},
		{
			name:           "service failure",
			query:          "",
			accountService: mock.NewAccountService(errors.New("forced")),
			expectedStatus: 500,
			expectedCalls:  1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := http.Config{
				Logger: zap.NewNop(),
				Services: http.Services{
					AccountService: c.accountService,
				},
			}

			s := httptest.NewServer(http.BuildRoutes(cfg))
			defer s.Close()
			url := fmt.Sprintf("%s/account?%s", s.URL, c.query)
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

			if c.accountService.ListAccountsCalled != c.expectedCalls {
				t.Fatalf("expected FetchAccount account to be called %d times, but it was called %d times", c.expectedCalls, c.accountService.CreateAccountCalled)
			}
		})
	}
}
