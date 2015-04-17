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
				&param{
					pType:        path_value_param_type,
					name:         "filepath",
					defaultValue: "index.html",
				},
			},
		},
	},
	{
		line: `post     /test/:ab/:cd                       Test.At(ab, cd=1, ef, gg=true)`,
		r: &routerLine{
			httpMethod: "POST",
			path:       "/test/:ab/:cd",
			structName: "Test",
			methodName: "At",
			params: []*param{
				&param{
					pType: path_param_type,
					name:  "ab",
				},
				&param{
					pType:        path_value_param_type,
					name:         "cd",
					defaultValue: "1",
				},
				&param{
					pType: query_string_param_type,
					name:  "ef",
				},
				&param{
					pType:        query_string_value_param_type,
					name:         "ef",
					defaultValue: "true",
				},
			},
		},
	},
	{
		line: `Delete     /test/:ab                       Test.Special(ab, cd= " , ", ef=  true for special test  , gg)`,
		r: &routerLine{
			httpMethod: "DELETE",
			path:       "/test/:ab/:cd",
			structName: "Test",
			methodName: "Special",
			params: []*param{
				&param{
					pType: path_param_type,
					name:  "ab",
				},
				&param{
					pType:        path_value_param_type,
					name:         "cd",
					defaultValue: " , ",
				},
				&param{
					pType:        query_string_value_param_type,
					name:         "ef",
					defaultValue: "true for special test",
				},
				&param{
					pType: query_string_param_type,
					name:  "ef",
				},
			},
		},
	},
	{
		line: `put     /test/:ab/:cd/:ef                       Test.Special(ab)`,
		r: &routerLine{
			httpMethod: "PUT",
			path:       "/test/:ab/:cd/:ef",
			structName: "Test",
			methodName: "Special",
			params: []*param{
				&param{
					pType: path_param_type,
					name:  "ab",
				},
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
