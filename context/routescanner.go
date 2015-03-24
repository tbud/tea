package context

import (
	"bytes"
	"fmt"
	. "github.com/tbud/x/builtin"
	"github.com/tbud/x/container/set"
	"strings"
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

type routeScanner struct {
	step func(*routeScanner, int) int

	// Error that happened, if any.
	err error

	data       []byte // store data load from file
	parseBuf   []byte // save parsed key or value
	bufType    int    //buf type
	bracketNum int    // save bracket num

	imports  map[string]set.StringSet
	includes []string
	routes   []string
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
		parseInclude(r)
	case bufInRoute:
		parseRoute(r)
	}

	r.parseBuf = r.parseBuf[:0]
	r.bufType = bufInUnknown
	r.step = stateBegin
	return stateBegin(r, c)
}

func parseImport(r *routeScanner) {
	var importList []string

	if bytes.Contains(r.parseBuf, []byte("(")) {

		buf := bytes.TrimPrefix(bytes.TrimSpace(r.parseBuf), []byte("import"))
		buf = bytes.TrimFunc(buf, func(r rune) bool {
			return r == '(' || r == ')'
		})
		importList = strings.Split(string(buf), "\n")
	} else {
		buf := bytes.TrimPrefix(bytes.TrimSpace(r.parseBuf), []byte("import"))
		importList = []string{string(buf)}
	}

	for _, importLine := range importList {
		importLine = strings.TrimSpace(importLine)
		if strings.ContainsAny(importLine, " \t") {

		} else {

		}
	}
}

func parseInclude(r *routeScanner) {

}

func parseRoute(r *routeScanner) {

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

func includeRoute(rootPath string, importAppPath string) (routers []router, err error) {
	defer Catch(func(ierr interface{}) {
		if errr, ok := ierr.(error); ok {
			err = errr
		}
		Log.Error("Catch error: %v", ierr)
	})

	return nil, nil
}
