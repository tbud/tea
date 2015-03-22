package context

import (
	// "strings"
	"testing"
)

var routerLines = []struct {
	num  int
	line string
	err  error
}{
	{1, "Get / index.html", nil},
	{3, "#Get / index.html", nil},
}

func TestAddRouter(t *testing.T) {
	pRouter := &Router{}
	for _, routerLine := range routerLines {
		if err := pRouter.AddRoute(routerLine.line, routerLine.num); err != nil {

		}
	}

}
