package context

import (
	"testing"
)

func TestIncludeRoute(t *testing.T) {
	if routes, err := includeRoute("/", "github.com/tbud/tea/archetype"); err != nil {
		t.Error(err)
	}
}
