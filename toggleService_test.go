package togglr_test

import (
	"context"
	"sync"
	"testing"

	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/mock"
	"github.com/togglr-io/togglr/rules"
	"github.com/togglr-io/togglr/uid"
	"go.uber.org/zap"
)

func Test_DefaultToggleService(t *testing.T) {
	// SETUP
	ms := mock.NewMetadataService(nil)
	mockTS := mock.NewToggleService(nil)
	ts := togglr.NewToggleService(mockTS, ms, zap.NewNop())

	rules := rules.Rules{
		{
			Expr: rules.Expression{
				Type: rules.ExprTypeBinary,
				Binary: rules.NewBinary(
					rules.NewIdent("test-key"),
					rules.NewString("test-value"),
					rules.BinOpEq,
				),
			},
		},
		{
			Expr: rules.Expression{
				Type: rules.ExprTypeBinary,
				Binary: rules.NewBinary(
					rules.NewIdent("another-key"),
					rules.NewString("another-value"),
					rules.BinOpEq,
				),
			},
		},
	}

	expectedKeys := []string{"test-key", "another-key"}

	toggle := togglr.Toggle{
		Key:   "test-toggle",
		Rules: rules,
	}

	wg := sync.WaitGroup{}
	ms.PushKeysFn = func(ctx context.Context, accountID uid.UID, keys ...string) error {
		defer wg.Done()
		var found bool
		for _, expected := range expectedKeys {
			found = false
			for _, key := range keys {
				if key == expected {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf("expected key %s", expected)
			}
		}

		return nil
	}

	// RUN
	wg.Add(2)
	_, err := ts.CreateToggle(context.TODO(), toggle)
	if err != nil {
		t.Fatalf("failed to create toggle: %s", err)
	}
	err = ts.UpdateToggle(context.TODO(), togglr.UpdateToggleReq{
		Rules: rules,
	})
	if err != nil {
		t.Fatalf("failed to update toggle: %s", err)
	}
	wg.Wait()

	if ms.PushKeysCalled != 2 {
		t.Fatalf("expected MetadataService.PushKeys to be called 1 time, not %d", ms.PushKeysCalled)
	}

	if mockTS.CreateToggleCalled != 1 {
		t.Fatalf("expected ToggleService.CreateTogggle to be called 1 time, not %d", mockTS.CreateToggleCalled)
	}

	if mockTS.UpdateToggleCalled != 1 {
		t.Fatalf("expected ToggleService.UpdateTogggle to be called 1 time, not %d", mockTS.UpdateToggleCalled)
	}
}
