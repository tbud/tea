package router

import (
	"reflect"
	"testing"
)

var testRouteLines = []struct {
	line string
	r    *routerLine
}{
	{
		line: `GET     /public/*filepath                       Assets.At("public", filepath = "index.html")`,
		r: &routerLine{
			httpMethod: "GET",
			path:       "/public/*filepath",
			structName: "Assets",
			methodName: "At",
			params: []*param{
				&param{
					pType:        fixed_value_type,
					defaultValue: "public",
				},
				// param{pType: default_value_type, name: "filepath", defaultValue: "index.html"},
			},
		},
	},
}

func TestParseRouteLine(t *testing.T) {
	for _, trl := range testRouteLines {
		if rl, err := parseRouterLine(trl.line); err != nil {
			t.Error(err)
			return
		} else if !reflect.DeepEqual(trl.r, rl) {
			t.Errorf("Want %v, Got %v", trl.r, rl)
		}
	}

}
