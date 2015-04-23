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
					name:         "gg",
					defaultValue: "true",
				},
			},
		},
	},
	{
		line: `Delete     /test/:ab                       Test.Special(ab, cd   = " , ", ef=  true for special test  , gg)`,
		r: &routerLine{
			httpMethod: "DELETE",
			path:       "/test/:ab",
			structName: "Test",
			methodName: "Special",
			params: []*param{
				&param{
					pType: path_param_type,
					name:  "ab",
				},
				&param{
					pType:        query_string_value_param_type,
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
					name:  "gg",
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

// func (r *routerLine) String() string {
// 	return fmt.Sprintf("%v\n", r.params)
// }

// func (p *param) String() string {
// 	return fmt.Sprintf("(%s, %s, %d, %v)", p.defaultValue, p.name, p.pType, p.typeExpr)
// }

func TestParseRouteLine(t *testing.T) {
	for _, trl := range testRouteLines {
		if rl, err := parseRouterLine(trl.line); err != nil {
			t.Error(err)
			return
		} else {
			rl.pathParams = nil
			if !reflect.DeepEqual(trl.r, rl) {
				t.Errorf("Want %v, Got %v", trl.r, rl)
			}
		}
	}
}

var testParamLines = []struct {
	paramLine string
	param     *param
	err       error
}{
	{
		paramLine: `ab`,
		param: &param{
			pType: path_param_type,
			name:  "ab",
		},
		err: nil,
	},
	{
		paramLine: `ab=1`,
		param: &param{
			pType:        path_value_param_type,
			name:         "ab",
			defaultValue: "1",
		},
		err: nil,
	},
	{
		paramLine: `cd`,
		param: &param{
			pType: query_string_param_type,
			name:  "cd",
		},
		err: nil,
	},
	{
		paramLine: `cd="abc"`,
		param: &param{
			pType:        query_string_value_param_type,
			name:         "cd",
			defaultValue: "abc",
		},
		err: nil,
	},
	{
		paramLine: `false`,
		param: &param{
			pType:        fixed_value_type,
			defaultValue: "false",
		},
		err: nil,
	},
	{
		paramLine: `"cd1"`,
		param: &param{
			pType:        fixed_value_type,
			defaultValue: "cd1",
		},
		err: nil,
	},
}

func TestRouterActionParam(t *testing.T) {
	for _, pl := range testParamLines {
		r := &routerLine{
			pathParams: []*param{
				&param{name: "ab"},
			},
		}

		if err := parseRouterActionParam(pl.paramLine, r); err != pl.err {
			t.Errorf("Parse param err: %v, param line: %s", err, pl.paramLine)
		} else {
			if !reflect.DeepEqual(r.params[0], pl.param) {
				t.Errorf("Parse line: %s, want: %v, got: %v", pl.paramLine, pl.param, r.params[0])
			}
		}
	}
}
