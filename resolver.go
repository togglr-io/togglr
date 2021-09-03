package togglr

import (
	"context"

	"github.com/togglr-io/togglr/rules"
	"github.com/togglr-io/togglr/uid"
)

type DefaultResolver struct {
	ts ToggleService
}

func NewResolver(ts ToggleService) DefaultResolver {
	return DefaultResolver{
		ts: ts,
	}
}

func (r DefaultResolver) Resolve(ctx context.Context, accountID uid.UID, md rules.Metadata) (ResolvedToggles, error) {
	resolved := make(ResolvedToggles)
	toggles, err := r.ts.ListToggles(ctx, ListTogglesReq{AccountID: accountID})
	if err != nil {
		return nil, err
	}

	for _, toggle := range toggles {
		resolved[toggle.Key] = rules.EvaluateRules(md, toggle.Rules...)
	}

	return resolved, nil
}
