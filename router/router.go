package context

import (
	. "github.com/tbud/x/builtin"
	// "io/ioutil"
	// "path/filepath"
	// "regexp"
	// "strings"
)

type router struct {
	httpMethod string
	path       string
	prefix     string
	structName string
	action     string
	params     []param
}

type param struct {
	name         string
	pType        paramType
	defaultValue interface{}
}

type paramType uint8

const (
	default_type paramType = iota
	fix_value_type
	default_value_type
)

func includeRoute(rootPath string, importAppPath string) (routers []router, err error) {
	defer Catch(func(ierr interface{}) {
		if errr, ok := ierr.(error); ok {
			err = errr
		}
		Log.Error("Catch error: %v", ierr)
	})

	return nil, nil
}

// func LoadRouter(file string) (pRouter *Router, err error) {
// 	if !filepath.IsAbs(file) {
// 		if file, err = filepath.Abs(file); err != nil {
// 			return nil, err
// 		}
// 	}

// 	pRouter = &Router{}

// 	var fileBuf []byte
// 	if fileBuf, err = ioutil.ReadFile(file); err != nil {
// 		for num, line := range strings.Split(string(fileBuf), "\n") {
// 			line = strings.TrimSpace(line)
// 			if len(line) == 0 || line[0] == '#' {
// 				continue
// 			}

// 			if err = pRouter.AddRoute(line, num); err != nil {
// 				return nil, err
// 			}
// 		}

// 		return pRouter, nil
// 	}

// 	return nil, err
// }

// func (r *Router) AddRoute(line string, num int) (err error) {
// 	return nil
// }

// var routePattern *regexp.Regexp = regexp.MustCompile("(?i)^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD|WS)[ \t]+([^ \t]+)[ \t]+(.+)$")

// func parseRouteLine(line string) (method, path, action string, found bool) {
// 	var matches []string = routePattern.FindStringSubmatch(line)
// 	if matches == nil {
// 		return
// 	}

// 	return matches[1], matches[2], strings.TrimSpace(matches[3]), true
// }

// func parseImportLine(line string) {

// }
