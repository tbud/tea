package context

import (
	"strings"
	"testing"
)

var routerLines = `

`

func TestAddRouter(t *testing.T) {
	pRouter := &Router{}
	for num, line := range strings.Split(routerLines, "\n") {
		if err = pRouter.AddRoute(line, num); err != nil {

		}
	}

}
