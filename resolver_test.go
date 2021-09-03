package togglr_test

import (
	"testing"

	"github.com/togglr-io/togglr"
)

func Test_DefaultResolver(t *testing.T) {
	// SETUP
	ts := mock.NewToggleService(nil)
	resolver := togglr.NewResolver(ts)
	metadata := togglr.Metadata{
		"userType": "admin",
		""
	}
}
