package router

import (
	"testing"
)

func TestIncludeRoute(t *testing.T) {
	if _, err := includeRoute("/", "github.com/tbud/tea/archetype"); err != nil {
		t.Error(err)
	}
}
