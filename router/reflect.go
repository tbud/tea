package context

import (
	. "github.com/tbud/tea/context"
	. "github.com/tbud/x/builtin"
	// "io/ioutil"
	"github.com/tbud/x/container/set"
	"go/ast"
	"go/build"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	// "regexp"
	"fmt"
	"strings"
)

type router struct {
	httpMethod string
	path       string
	prefix     string
	structName string
	actions    []action
}

type action struct {
	name   string
	params []param
}

type param struct {
	name         string
	rType        reflect.Type
	pType        paramType
	defaultValue interface{}
}

type paramType uint8

const (
	default_type paramType = iota
	fix_value_type
	default_value_type
)

type importStruct struct {
	structName string
	actions    []action
}

const (
	route_file_path     = "conf/routes"
	tea_controller_path = "app/controllers"
)

func parseRouteScanner(r *routeScanner, rootPath string) (routers []router, err error) {
	if routers, err = parseRouteScannerRouter(r); err != nil {
		Log.Error("%v", err)
		return
	}

	if len(r.includes) > 0 {
		for rootPath, importAppPath := range r.includes {
			if r, errr := includeRoute(rootPath, importAppPath); errr != nil {
				Log.Error("include router (%s, %s) error: %v", rootPath, importAppPath, errr)
				return
			} else {
				routers = append(routers, r...)
			}
		}
	}

	return routers, nil
}

func parseRouteScannerRouter(r *routeScanner) (routers []router, err error) {

	return nil, nil
}

func parseDirController(rootImportPath string) (importStructs []importStruct, err error) {
	var pkg *build.Package
	if pkg, err = build.ImportDir(rootImportPath, build.FindOnly); err != nil {
		return err
	}

	root := pkg.Dir
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Log.Error("Error scanning app source: %v", err)
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		pkgImportPath := filepath.Join(rootImportPath, path[len(root):])

		var pkgs map[string]*ast.Package
		fset := token.NewFileSet()
		pkgs, err = parser.ParseDir(fset, path, func(f os.FileInfo) bool {
			return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && filepath.Ext(f.name) == ".go"
		}, 0)

		if err != nil {
			return err
		}

		// Skip "main" packages.
		delete(pkgs, "main")

		// If there is no code in this directory, skip it.
		if len(pkgs) == 0 {
			return nil
		}

		// There should be only one package in this directory.
		if len(pkgs) > 1 {
			Log.Warn("Most unexpected! Multiple packages in a single directory: %v", pkgs)
		}

		var pkg *ast.Package
		for _, v := range pkgs {
			pkg = v
		}

		importStructs = append(importStructs, processPackage(fset, pkgImportPath, path, pkg)...)
	})

	return
}

func processPackage(fset *token.FileSet, pkgImportPath, pkgPath string, pkg *ast.Package) []importStruct {
	fmt.Printf("%#v\n", pkg)
	// for _, file := range pkg.Files {
	// 	imports := map[string]string{}

	// 	for _, decl := range file.Decls {

	// 	}
	// }
	return nil
}

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
