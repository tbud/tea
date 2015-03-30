package context

import (
	. "github.com/tbud/tea/context"
	// "io/ioutil"
	"go/ast"
	"go/build"
	"go/parser"
	// "go/scanner"
	"go/token"
	"os"
	"path/filepath"
	// "regexp"
	"fmt"
	"strings"
)

type router struct {
	httpMethod string
	path       string
	prefix     string
	structName string
	methods    []method
}

type method struct {
	name   string
	params []*param
}

type param struct {
	name         string
	typeExpr     TypeExpr
	pType        paramType
	defaultValue interface{}
}

type paramType uint8

const (
	default_type paramType = iota
	fix_value_type
	default_value_type
)

type controller struct {
	structName  string
	importPath  string
	packageName string
	methods     []*method
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

	return routers, nil
}

func parseRouteScannerRouter(r *routeScanner) (routers []router, err error) {

	return nil, nil
}

func parseDirController(rootImportPath string) (controllers []controller, err error) {
	var pkg *build.Package
	rootImportPath = filepath.Join(rootImportPath, tea_controller_path)
	if pkg, err = build.Import(rootImportPath, "", build.FindOnly); err != nil {
		Log.Error("%v", err)
		return nil, err
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
			return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && filepath.Ext(f.Name()) == ".go"
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

		controllers = append(controllers, processPackage(fset, pkgImportPath, path, pkg)...)
		return nil
	})

	return
}

func processPackage(fset *token.FileSet, pkgImportPath, pkgPath string, pkg *ast.Package) []controller {

	for _, file := range pkg.Files {
		// imports := map[string]string{}

		for _, decl := range file.Decls {
			appendAction(fset, decl, pkgImportPath, pkg.Name)
		}
	}
	return nil
}

func appendStruct(pkgImportPath string, pkg *ast.Package, decl ast.Decl, fset *token.FileSet) []*controller {
	spec, found := getStructTypeDecl(decl, fset)
	if !found {
		return nil
	}

	// structType := spec.Type.(*ast.StructType)

	imStruct := &controller{
		structName:  spec.Name.Name,
		importPath:  pkgImportPath,
		packageName: pkg.Name,
	}

	return []*controller{imStruct}
}

func appendAction(fset *token.FileSet, decl ast.Decl, pkgImportPath, pkgName string) {
	// Func declaration?
	funcDecl, ok := decl.(*ast.FuncDecl)
	if !ok {
		return
	}

	// Have a receive ? is receive is nil, when func is not a struct func
	if funcDecl.Recv == nil {
		return
	}

	// Is it public?
	if !funcDecl.Name.IsExported() {
		return
	}

	// Does it return a result?
	if funcDecl.Type.Results != nil {
		return
	}

	m := &method{name: funcDecl.Name.Name}

	for _, field := range funcDecl.Type.Params.List {
		for _, name := range field.Names {
			// var importPath string
			typeExpr := newTypeExpr(pkgName, field.Type)
			if !typeExpr.valid {
				return // we didn't understand one of the args. Ignore this action.
			}

			m.params = append(m.params, &param{name: name.Name, typeExpr: typeExpr})
		}
	}

	fmt.Printf("%#v\n", m)

	ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		return true
	})
}

type TypeExpr struct {
	expr     string
	pkgName  string
	pkgIndex int
	valid    bool
}

var _BUILTIN_TYPES = map[string]bool{
	"bool":       true,
	"byte":       true,
	"complex128": true,
	"complex64":  true,
	"error":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"int8":       true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uint8":      true,
	"uintptr":    true,
}

func isBuiltinType(name string) (found bool) {
	_, found = _BUILTIN_TYPES[name]
	return
}

func newTypeExpr(pkgName string, expr ast.Expr) TypeExpr {
	switch t := expr.(type) {
	case *ast.Ident:
		if isBuiltinType(t.Name) {
			pkgName = ""
		}
		return TypeExpr{t.Name, pkgName, 0, true}
	case *ast.SelectorExpr:
		e := newTypeExpr(pkgName, t.X)
		return TypeExpr{t.Sel.Name, e.expr, 0, e.valid}
	case *ast.StarExpr:
		e := newTypeExpr(pkgName, t.X)
		return TypeExpr{"*" + e.expr, e.pkgName, e.pkgIndex + 1, e.valid}
	case *ast.ArrayType:
		e := newTypeExpr(pkgName, t.Elt)
		return TypeExpr{"[]" + e.expr, e.pkgName, e.pkgIndex + 2, e.valid}
	case *ast.Ellipsis:
		e := newTypeExpr(pkgName, t.Elt)
		return TypeExpr{"[]" + e.expr, e.pkgName, e.pkgIndex + 3, e.valid}
	default:
		Log.Warn("Failed to generate name for field. Make sure the field name is valid.")
	}
	return TypeExpr{valid: false}
}

func getStructTypeDecl(decl ast.Decl, fset *token.FileSet) (spec *ast.TypeSpec, found bool) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return
	}

	if genDecl.Tok != token.TYPE {
		return
	}

	if len(genDecl.Specs) == 0 {
		pos := fset.Position(decl.Pos())
		Log.Warn("Surprising: %s:%d Decl contains no specifications", pos.Filename, pos.Line)
		return
	}

	spec = genDecl.Specs[0].(*ast.TypeSpec)
	_, found = spec.Type.(*ast.StructType)

	return
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
