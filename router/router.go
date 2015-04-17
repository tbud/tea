package router

import (
	"fmt"
	"strings"

	"regexp"
)

var routePattern *regexp.Regexp = regexp.MustCompile("(?i)^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD|WS)[ \t]+([^ \t]+)[ \t]+(.+)$")

func parseRouterLine(line string) (r *routerLine, err error) {
	matches := routePattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("Line couldn't match reg: %s", line)
	}

	pathLine := matches[2]
	r = &routerLine{
		httpMethod: strings.ToUpper(matches[1]),
		path:       pathLine,
	}

	if err = parseRouterPath(pathLine, r); err != nil {
		return nil, err
	}

	if err = parseRouterAction(matches[3], r); err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", r)
	for _, p := range r.params {
		fmt.Printf("%v\n", p)
	}

	return r, nil
}

func parseRouterPath(pathLine string, r *routerLine) error {
	pathLine = strings.TrimSpace(pathLine)

	var (
		getParamName = false
		nameBuf      []rune
	)

	for _, c := range pathLine {
		switch c {
		default:
			if getParamName {
				nameBuf = append(nameBuf, c)
			}
		case ':', '*':
			getParamName = true
		case '/':
			if getParamName {
				r.params = append(r.params, &param{pType: path_param_type, name: string(nameBuf)})
				nameBuf = nameBuf[:0]
			}
		}
	}

	if getParamName && len(nameBuf) > 0 {
		r.params = append(r.params, &param{pType: path_param_type, name: string(nameBuf)})
	}

	return nil
}

func parseRouterAction(actionLine string, r *routerLine) error {
	actionLine = strings.TrimSpace(actionLine)

	return nil
}
