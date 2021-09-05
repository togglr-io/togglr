package togglr_test

import (
	"context"
	"testing"

	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/mock"
	"github.com/togglr-io/togglr/rules"
	"github.com/togglr-io/togglr/uid"
)

func listToggles(ctx context.Context, req togglr.ListTogglesReq) ([]togglr.Toggle, error) {
	return []togglr.Toggle{
		{
			ID:  uid.New(),
			Key: "admin-feature",
			Rules: rules.Rules{
				{
					Op: rules.BinOpAnd,
					Expr: rules.Expression{
						Type: rules.ExprTypeBinary,
						Binary: rules.NewBinary(
							rules.NewIdent("userType"),
							rules.NewString("admin"),
							rules.BinOpEq,
						),
					},
				},
				{
					Op: rules.BinOpAnd,
					Expr: rules.Expression{
						Type: rules.ExprTypeBinary,
						Binary: rules.NewBinary(
							rules.NewIdent("hasFlag"),
							rules.NewBool(true),
							rules.BinOpEq,
						),
					},
				},
			},
		},
		{
			ID:  uid.New(),
			Key: "user-feature",
			Rules: rules.Rules{
				{
					Op: rules.BinOpAnd,
					Expr: rules.Expression{
						Type: rules.ExprTypeBinary,
						Binary: rules.NewBinary(
							rules.NewIdent("userType"),
							rules.NewString("admin"),
							rules.BinOpNotEq,
						),
					},
				},
				{
					Op: rules.BinOpOr,
					Expr: rules.Expression{
						Type: rules.ExprTypeBinary,
						Binary: rules.NewBinary(
							rules.NewIdent("iDontExist"),
							rules.NewString("whatever"),
							rules.BinOpEq,
						),
					},
				},
			},
		},
	}, nil
}

func Test_DefaultResolver(t *testing.T) {
	// SETUP
	ctx := context.TODO()
	ts := mock.NewToggleService(nil)
	ts.ListTogglesFn = listToggles
	resolver := togglr.NewResolver(ts)
	metadata := rules.Metadata{
		"userType": rules.NewString("admin"),
		"hasFlag":  rules.NewBool(true),
	}

	// RUN
	resolved, err := resolver.Resolve(ctx, uid.New(), metadata)
	if err != nil {
		t.Fatalf("failed to resolve toggles: %s", err)
	}

	admin, ok := resolved["admin-feature"]
	if !ok {
		t.Fatalf("expected 'admin-feature' flag to be present")
	}

	if !admin {
		t.Fatalf("expected 'admin-feature' flag to be true")
	}

	user, ok := resolved["user-feature"]
	if !ok {
		t.Fatalf("expected 'user-feature' flag to be present")
	}

	if user {
		t.Fatalf("expected 'user-feature' flag to be false")
	}
}
