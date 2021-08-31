package togglr

import (
	"context"

	"github.com/togglr-io/togglr/rules"
	"github.com/togglr-io/togglr/uid"
	"go.uber.org/zap"
)

// A DefaultToggleService provides a default implementation of the ToggleService interface that
// wraps another ToggleService and performs some additional business logic
type DefaultToggleService struct {
	ts ToggleService
	ms MetadataService

	log *zap.Logger
}

// NewToggleService returns a new DeafultToggleService.
func NewToggleService(ts ToggleService, ms MetadataService, logger *zap.Logger) DefaultToggleService {
	return DefaultToggleService{
		ts:  ts,
		ms:  ms,
		log: logger,
	}
}

// TODO (etate): This is very specifically checking binary and ident expressions only. Might need to
// generalize once things like Unary exprs exist
func extractKeys(expr rules.Expr) []string {
	keys := []string{}
	switch v := expr.(type) {
	case rules.Expression:
		switch v.Type {
		case rules.ExprTypeBinary:
			keys = append(extractKeys(v.Binary))
		case rules.ExprTypeIdent:
			keys = append(extractKeys(v.Ident))
		}
	case rules.Binary:
		keys = append(keys, extractKeys(v.Left)...)
		keys = append(keys, extractKeys(v.Right)...)
	case rules.Ident:
		keys = append(keys, v.Value)
	}

	return keys
}

func (s DefaultToggleService) pushKeys(ctx context.Context, accountID uid.UID, rules rules.Rules) {
	defer s.log.Sync()
	if rules == nil || len(rules) == 0 {
		return
	}

	// collect keys from Rules
	keys := []string{}
	for _, rule := range rules {
		keys = append(keys, extractKeys(rule.Expr)...)
	}

	if err := s.ms.PushKeys(ctx, accountID, keys...); err != nil {
		s.log.Error("failed to push metadata keys", zap.Error(err))
	}
}

func (s DefaultToggleService) CreateToggle(ctx context.Context, toggle Toggle) (uid.UID, error) {
	// push keys asynchronously so we don't keep the caller waiting
	go s.pushKeys(ctx, toggle.AccountID, toggle.Rules)

	return s.ts.CreateToggle(ctx, toggle)
}

func (s DefaultToggleService) UpdateToggle(ctx context.Context, req UpdateToggleReq) error {
	go s.pushKeys(ctx, req.AccountID, req.Rules)

	return s.ts.UpdateToggle(ctx, req)
}

func (s DefaultToggleService) FetchToggle(ctx context.Context, id uid.UID) (Toggle, error) {
	return s.ts.FetchToggle(ctx, id)
}

func (s DefaultToggleService) ListToggles(ctx context.Context, req ListTogglesReq) ([]Toggle, error) {
	return s.ts.ListToggles(ctx, req)
}

func (s DefaultToggleService) DeleteToggle(ctx context.Context, id uid.UID) error {
	return s.ts.DeleteToggle(ctx, id)
}
