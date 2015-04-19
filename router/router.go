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
	for _, p := range r.pathParams {
		fmt.Printf("%v\n", p)
	}
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
				r.pathParams = append(r.pathParams, &param{pType: path_param_type, name: string(nameBuf)})
				nameBuf = nameBuf[:0]
			}
		}
	}

	if getParamName && len(nameBuf) > 0 {
		r.pathParams = append(r.pathParams, &param{pType: path_param_type, name: string(nameBuf)})
	}

	return nil
}

func parseRouterAction(actionLine string, r *routerLine) (err error) {
	actionLine = strings.TrimSpace(actionLine)

	var (
		bInQuote     = false
		bInStringEsc = false
		buf          []rune
	)

	for _, c := range actionLine {
		switch c {
		case '"':
			if !bInStringEsc {
				bInQuote = !bInQuote
			} else {
				bInStringEsc = false
			}
			buf = append(buf, c)
			continue
		case '\\':
			bInStringEsc = !bInStringEsc
			buf = append(buf, c)
			continue
		case 'b', 'f', 'n', 'r', 't', '/', 'u':
			if bInStringEsc {
				buf = append(buf, c)
				bInStringEsc = false
				continue
			}
		}

		if bInQuote {
			buf = append(buf, c)
		} else {
			switch c {
			case '.':
				if len(r.structName) == 0 {
					r.structName = string(buf)
					buf = buf[:0]
				} else {
					return fmt.Errorf("find second '.' in route action: %s", actionLine)
				}
			case '(':
				if len(r.methodName) == 0 {
					r.methodName = string(buf)
					buf = buf[:0]
				} else {
					return fmt.Errorf("find second '(' in route action", actionLine)
				}
			case ',', ')':
				paramLine := string(buf)
				buf = buf[:0]
				if err = parseRouterActionParam(paramLine, r); err != nil {
					return err
				}
			default:
				buf = append(buf, c)
			}
		}
	}

	if len(buf) > 0 {
		return fmt.Errorf("Endless parse: '%s' in actionline: %s ", string(buf), actionLine)
	}

	return nil
}

func parseRouterActionParam(paramLine string, r *routerLine) error {
	paramLine = strings.TrimSpace(paramLine)
	quoteIndex := strings.Index(paramLine, "\"")
	equalIndex := strings.Index(paramLine, "=")

	if equalIndex > 0 {
		if quoteIndex > 0 {
			if quoteIndex > equalIndex {
				paramName := strings.TrimSpace(paramLine[0:equalIndex])
				paramValue := strings.TrimSpace(paramLine[equalIndex+1:])

			} else {

			}
		} else {

		}
	} else {

	}
	return nil
}

func (rl *routerLine) findPathParamByName(name string) *param {
	for _, p := range rl.pathParams {
		if p.name == name {
			return p
		}
	}

	return nil
}
