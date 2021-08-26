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

func Test_HandleTogglePost(t *testing.T) {
	id := uid.New().String()
	cases := []struct {
		name                string
		payload             string
		toggleService       *mock.ToggleService
		expectedStatus      int
		expectedCreateCalls int
		expectedUpdateCalls int
	}{
		{
			name:                "successful create",
			payload:             `{"key": "test-toggle"}`,
			toggleService:       mock.NewToggleService(nil),
			expectedStatus:      200,
			expectedCreateCalls: 1,
			expectedUpdateCalls: 0,
		},
		{
			name:                "bad request",
			payload:             `{"key": invalid"test-toggle"}`,
			toggleService:       mock.NewToggleService(nil),
			expectedStatus:      400,
			expectedCreateCalls: 0,
			expectedUpdateCalls: 0,
		},
		{
			name:                "failed create",
			payload:             `{"key": "test-toggle"}`,
			toggleService:       mock.NewToggleService(errors.New("forced")),
			expectedStatus:      500,
			expectedCreateCalls: 1,
			expectedUpdateCalls: 0,
		},
		{
			name:                "successful update",
			payload:             fmt.Sprintf(`{"id": "%s", "description": "New description"}`, id),
			toggleService:       mock.NewToggleService(nil),
			expectedStatus:      200,
			expectedCreateCalls: 0,
			expectedUpdateCalls: 1,
		},
		{
			name:                "failed update",
			payload:             fmt.Sprintf(`{"id": "%s", "description": "New description"}`, id),
			toggleService:       mock.NewToggleService(errors.New("forced")),
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
					ToggleService: c.toggleService,
				},
			}

			s := httptest.NewServer(http.BuildRoutes(cfg))
			defer s.Close()
			url := fmt.Sprintf("%s/toggle", s.URL)
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

			if c.toggleService.CreateToggleCalled != c.expectedCreateCalls {
				t.Fatalf("expected CreateToggle toggle to be called %d times, but it was called %d times", c.expectedCreateCalls, c.toggleService.CreateToggleCalled)
			}

			if c.toggleService.UpdateToggleCalled != c.expectedUpdateCalls {
				t.Fatalf("expected UpdateToggle toggle to be called %d times, but it was called %d times", c.expectedUpdateCalls, c.toggleService.UpdateToggleCalled)
			}
		})
	}
}

func Test_HandleToggleGet(t *testing.T) {
	cases := []struct {
		name           string
		query          string
		toggleService  *mock.ToggleService
		expectedStatus int
		expectedCalls  int
	}{
		{
			name:           "successful test",
			query:          "",
			toggleService:  mock.NewToggleService(nil),
			expectedStatus: 200,
			expectedCalls:  1,
		},
		{
			name:           "service failure",
			query:          "",
			toggleService:  mock.NewToggleService(errors.New("forced")),
			expectedStatus: 500,
			expectedCalls:  1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := http.Config{
				Logger: zap.NewNop(),
				Services: http.Services{
					ToggleService: c.toggleService,
				},
			}

			s := httptest.NewServer(http.BuildRoutes(cfg))
			defer s.Close()
			url := fmt.Sprintf("%s/toggle?%s", s.URL, c.query)
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

			if c.toggleService.ListTogglesCalled != c.expectedCalls {
				t.Fatalf("expected FetchToggle toggle to be called %d times, but it was called %d times", c.expectedCalls, c.toggleService.CreateToggleCalled)
			}
		})
	}
}

func Test_HandleToggleDelete(t *testing.T) {
	cases := []struct {
		name           string
		id             string
		toggleService  *mock.ToggleService
		expectedStatus int
		expectedCalls  int
	}{
		{
			name:           "successful test",
			id:             "c149f08b-b0fa-4a5d-8a6c-03ac992aa454",
			toggleService:  mock.NewToggleService(nil),
			expectedStatus: 204,
			expectedCalls:  1,
		},
		{
			name:           "bad request",
			id:             "123",
			toggleService:  mock.NewToggleService(nil),
			expectedStatus: 400,
			expectedCalls:  0,
		},
		{
			name:           "service failure",
			id:             "c149f08b-b0fa-4a5d-8a6c-03ac992aa454",
			toggleService:  mock.NewToggleService(errors.New("forced")),
			expectedStatus: 500,
			expectedCalls:  1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := http.Config{
				Logger: zap.NewNop(),
				Services: http.Services{
					ToggleService: c.toggleService,
				},
			}

			s := httptest.NewServer(http.BuildRoutes(cfg))
			defer s.Close()
			url := fmt.Sprintf("%s/toggle/%s", s.URL, c.id)
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

			if c.toggleService.DeleteToggleCalled != c.expectedCalls {
				t.Fatalf("expected DeleteToggle to be called %d times, but it was called %d times", c.expectedCalls, c.toggleService.CreateToggleCalled)
			}
		})
	}
}
