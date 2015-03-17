package context

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Router struct {
	routers []*Router
}

func LoadRouter(file string) (pRouter *Router, err error) {
	if !filepath.IsAbs(file) {
		file = filepath.Abs(file)
	}

	pRouter = &Router{}

	var fineBuf []byte
	if fileBuf, err = ioutil.ReadFile(file); err != nil {
		for num, line := range strings.Split(string(fileBuf), "\n") {
			line = strings.TrimSpace(line)
			if len(line) == 0 || line[0] == '#' {
				continue
			}

			if err = pRouter.AddRoute(line, num); err != nil {
				return nil, err
			}
		}

		return pRouter, nil
	}

	return nil, err
}

func (r *Router) AddRoute(line string, num int) (err error) {

}
