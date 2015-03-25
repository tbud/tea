package context

import (
	"bytes"
	"fmt"
	"github.com/tbud/x/container/set"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	bufInUnknown = iota
	bufInImport
	bufInInclude
	bufInRoute
)

const (
	scanContinue  = iota // Continue, uninteresting byte
	scanAppendBuf        // byte need to append buf
	scanError            // hit an error, scanner.err.
)

type routeScanner struct {
	step func(*routeScanner, int) int

	// Error that happened, if any.
	err error

	data       []byte // store data load from file
	parseBuf   []byte // save parsed key or value
	bufType    int    //buf type
	bracketNum int    // save bracket num

	imports  map[string]*set.StringSet
	includes map[string]string
	routes   []string
}

func (r *routeScanner) open(fileName string) (err error) {
	r.init()

	if !filepath.IsAbs(fileName) {
		return fmt.Errorf("file '%s' is not absolute path", fileName)
	}

	r.data, err = ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	for _, c := range r.data {
		switch r.step(r, int(c)) {
		case scanError:
			return r.err
		case scanAppendBuf:
			r.parseBuf = append(r.parseBuf, c)
		}
	}

	if r.step(r, '\n') == scanError {
		return r.err
	}
	return nil
}

func (r *routeScanner) init() {
	r.step = stateBegin
	r.err = nil
	r.bufType = bufInUnknown
}

func stateBegin(r *routeScanner, c int) int {
	if c <= ' ' && isSpace(rune(c)) {
		return scanContinue
	}

	switch c {
	case '#':
		r.step = stateComment
		return scanContinue
	}

	r.step = stateParseLine
	return stateParseLine(r, c)
}

func stateParseLine(r *routeScanner, c int) int {
	if r.bufType == bufInUnknown {
		switch c {
		case ' ', '\t':
			keywords.checkKeyword(r, r.parseBuf)
			return scanAppendBuf
		case '(':
			r.bracketNum += 1
			keywords.checkKeyword(r, r.parseBuf)
			return scanAppendBuf
		}
	} else {
		switch c {
		case ')', '\r', '\n':
			if c == ')' {
				r.bracketNum -= 1
			}

			if r.bracketNum == 0 {
				r.step = stateEnd
				return scanAppendBuf
			}
			return scanAppendBuf
		}
	}

	return scanAppendBuf
}

func stateEnd(r *routeScanner, c int) int {
	switch r.bufType {
	case bufInUnknown:
		return r.error(c, "unkown route line: "+string(r.parseBuf))
	case bufInImport:
		parseImport(r)
	case bufInInclude:
		if parseInclude(r, c) == scanError {
			return scanError
		}
	case bufInRoute:
		parseRoute(r)
	}

	r.parseBuf = r.parseBuf[:0]
	r.bufType = bufInUnknown
	r.step = stateBegin
	return stateBegin(r, c)
}

func stateComment(r *routeScanner, c int) int {
	if c == '\n' || c == '\r' {
		r.step = stateBegin
		return scanContinue
	}
	return scanContinue
}

func stateError(r *routeScanner, c int) int {
	return scanError
}

func (r *routeScanner) error(c int, context string) int {
	r.step = stateError
	r.err = fmt.Errorf("invalid character '%c' , Error : %s", c, context)
	return scanError
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

/*********** parse buf content *************/
var importRegex = regexp.MustCompile("([^ \t]+)[ \t]+(.*)")

func getBlockList(r *routeScanner, prefix string) (list []string) {
	list = []string{}

	if bytes.Contains(r.parseBuf, []byte("(")) {
		buf := bytes.TrimPrefix(bytes.TrimSpace(r.parseBuf), []byte(prefix))
		buf = bytes.TrimFunc(buf, func(r rune) bool {
			return r == '(' || r == ')'
		})
		list = strings.Split(string(buf), "\n")
	} else {
		buf := bytes.TrimPrefix(bytes.TrimSpace(r.parseBuf), []byte(prefix))
		list = []string{string(buf)}
	}

	return
}

func parseImport(r *routeScanner) {
	for _, importLine := range getBlockList(r, "import") {
		importLine = strings.TrimSpace(importLine)
		var prefix, importPath string

		if strings.ContainsAny(importLine, " \t") {
			matches := importRegex.FindStringSubmatch(importLine)
			if matches == nil {
				continue
			}
			prefix, importPath = matches[1], matches[2]
			importPath = strings.Trim(strings.TrimSpace(importPath), "\"")
		} else {
			importPath = strings.Trim(importLine, "\"")
			prefix = filepath.Base(importPath)
		}

		var (
			m     *set.StringSet
			found bool
		)

		if m, found = r.imports[prefix]; !found {
			m = &set.StringSet{}
			r.imports[prefix] = m
		}
		m.Add(importPath)
	}
}

func parseInclude(r *routeScanner, c int) int {
	for _, includeLine := range getBlockList(r, "include") {
		includeLine = strings.TrimSpace(includeLine)

		if !strings.ContainsAny(includeLine, " \t") {
			return r.error(c, "parse include error: "+includeLine)
		}

		matches := importRegex.FindStringSubmatch(includeLine)
		if matches == nil {
			continue
		}

		r.includes[matches[1]] = matches[2]
	}

	return scanContinue
}

func parseRoute(r *routeScanner) {
	routePath := string(bytes.TrimSpace(r.parseBuf))
	r.routes = append(r.routes, routePath)
}

/************** keyword scanner *****************/
type keywordScanner struct {
	maxWordLen   int
	keywords     []string
	keywordsType []int
}

var keywords = keywordScanner{
	0,
	[]string{"import", "include"},
	[]int{bufInImport, bufInInclude},
}

func (k *keywordScanner) init() {
	for _, keyword := range k.keywords {
		wordLen := len(keyword)
		if wordLen > k.maxWordLen {
			k.maxWordLen = wordLen
		}
	}
}

func (k *keywordScanner) checkKeyword(r *routeScanner, buf []byte) {
	if len(buf) <= k.maxWordLen {
		value := string(buf)
		for i, keyword := range k.keywords {
			if value == keyword {
				r.bufType = k.keywordsType[i]
				return
			}
		}
	}

	r.bufType = bufInRoute
	return
}
