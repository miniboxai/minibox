package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/astrewrite"
	"github.com/kr/pretty"
	"golang.org/x/tools/imports"
)

const usage = `impl [-dir directory] <recv> <iface>

impl generates method stubs for recv to implement iface.


Examples:

impl 'f *File' io.Reader
impl Murmur hash.Hash
impl -dir $GOPATH/src/github.com/josharian/impl Murmur hash.Hash

Don't forget the single quotes around the receiver type
to prevent shell globbing.
`

var (
	flagSrcDir = flag.String("dir", "", "package source directory, useful for vendored code")
	fileInput  = flag.String("input", "i", "input source file")
)

// findInterface returns the import path and identifier of an interface.
// For example, given "http.ResponseWriter", findInterface returns
// "net/http", "ResponseWriter".
// If a fully qualified interface is given, such as "net/http.ResponseWriter",
// it simply parses the input.
func findInterface(iface string, srcDir string) (path string, id string, err error) {
	if len(strings.Fields(iface)) != 1 {
		return "", "", fmt.Errorf("couldn't parse interface: %s", iface)
	}

	srcPath := filepath.Join(srcDir, "__go_impl__.go")

	if slash := strings.LastIndex(iface, "/"); slash > -1 {
		// package path provided
		dot := strings.LastIndex(iface, ".")
		// make sure iface does not end with "/" (e.g. reject net/http/)
		if slash+1 == len(iface) {
			return "", "", fmt.Errorf("interface name cannot end with a '/' character: %s", iface)
		}
		// make sure iface does not end with "." (e.g. reject net/http.)
		if dot+1 == len(iface) {
			return "", "", fmt.Errorf("interface name cannot end with a '.' character: %s", iface)
		}
		// make sure iface has exactly one "." after "/" (e.g. reject net/http/httputil)
		if strings.Count(iface[slash:], ".") != 1 {
			return "", "", fmt.Errorf("invalid interface name: %s", iface)
		}
		return iface[:dot], iface[dot+1:], nil
	}

	src := []byte("package hack\n" + "var i " + iface)
	log.Printf("src: %s %s", src, srcDir)
	// If we couldn't determine the import path, goimports will
	// auto fix the import path.
	imp, err := imports.Process(srcPath, src, nil)
	if err != nil {
		return "", "", fmt.Errorf("couldn't parse interface: %s", iface)
	}

	// imp should now contain an appropriate import.
	// Parse out the import and the identifier.
	fset := token.NewFileSet()
	log.Printf("imp: %s", imp)
	f, err := parser.ParseFile(fset, srcPath, imp, 0)
	if err != nil {
		panic(err)
	}

	log.Printf("ParseFile: %# v", pretty.Formatter(f))
	if len(f.Imports) == 0 {
		return "", "", fmt.Errorf("unrecognized interface: %s", iface)
	}
	raw := f.Imports[0].Path.Value   // "io"
	path, err = strconv.Unquote(raw) // io
	if err != nil {
		panic(err)
	}
	decl := f.Decls[1].(*ast.GenDecl)      // var i io.Reader
	spec := decl.Specs[0].(*ast.ValueSpec) // i io.Reader
	sel := spec.Type.(*ast.SelectorExpr)   // io.Reader
	id = sel.Sel.Name                      // Reader
	return path, id, nil
}

// Pkg is a parsed build.Package.
type Pkg struct {
	*build.Package
	*token.FileSet
}

// typeSpec locates the *ast.TypeSpec for type id in the import path.
func typeSpec(path string, id string, srcDir string) (Pkg, *ast.TypeSpec, error) {
	pkg, err := build.Import(path, srcDir, 0)
	if err != nil {
		return Pkg{}, nil, fmt.Errorf("couldn't find package %s: %v", path, err)
	}

	fset := token.NewFileSet() // share one fset across the whole package
	for _, file := range pkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, 0)
		if err != nil {
			continue
		}

		for _, decl := range f.Decls {
			decl, ok := decl.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				continue
			}
			for _, spec := range decl.Specs {
				spec := spec.(*ast.TypeSpec)
				if spec.Name.Name != id {
					continue
				}
				return Pkg{Package: pkg, FileSet: fset}, spec, nil
			}
		}
	}
	return Pkg{}, nil, fmt.Errorf("type %s not found in %s", id, path)
}

// gofmt pretty-prints e.
func (p Pkg) gofmt(e ast.Expr) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, p.FileSet, e)
	return buf.String()
}

// fullType returns the fully qualified type of e.
// Examples, assuming package net/http:
// 	fullType(int) => "int"
// 	fullType(Handler) => "http.Handler"
// 	fullType(io.Reader) => "io.Reader"
// 	fullType(*Request) => "*http.Request"
func (p Pkg) fullType(e ast.Expr) string {
	ast.Inspect(e, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			// Using typeSpec instead of IsExported here would be
			// more accurate, but it'd be crazy expensive, and if
			// the type isn't exported, there's no point trying
			// to implement it anyway.
			if n.IsExported() {
				n.Name = p.Package.Name + "." + n.Name
			}
		case *ast.SelectorExpr:
			return false
		}
		return true
	})
	return p.gofmt(e)
}

func (p Pkg) params(field *ast.Field) []Param {
	var params []Param
	typ := p.fullType(field.Type)
	for _, name := range field.Names {
		params = append(params, Param{Name: name.Name, Type: typ})
	}
	// Handle anonymous params
	if len(params) == 0 {
		params = []Param{Param{Type: typ}}
	}
	return params
}

// Method represents a method signature.
type Method struct {
	Recv string
	Func
}

// Func represents a function signature.
type Func struct {
	Name   string
	Params []Param
	Res    []Param
}

// Param represents a parameter in a function or method signature.
type Param struct {
	Name string
	Type string
}

func (p Pkg) funcsig(f *ast.Field) Func {
	fn := Func{Name: f.Names[0].Name}
	typ := f.Type.(*ast.FuncType)
	if typ.Params != nil {
		for _, field := range typ.Params.List {
			fn.Params = append(fn.Params, p.params(field)...)
		}
	}
	if typ.Results != nil {
		for _, field := range typ.Results.List {
			fn.Res = append(fn.Res, p.params(field)...)
		}
	}
	return fn
}

// The error interface is built-in.
var errorInterface = []Func{{
	Name: "Error",
	Res:  []Param{{Type: "string"}},
}}

// funcs returns the set of methods required to implement iface.
// It is called funcs rather than methods because the
// function descriptions are functions; there is no receiver.
// func funcs(iface string, srcDir string) ([]Func, error) {
// 	// Special case for the built-in error interface.
// 	if iface == "error" {
// 		return errorInterface, nil
// 	}

// 	// Locate the interface.
// 	path, id, err := findInterface(iface, srcDir)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Parse the package and find the interface declaration.
// 	p, spec, err := typeSpec(path, id, srcDir)
// 	if err != nil {
// 		return nil, fmt.Errorf("interface %s not found: %s", iface, err)
// 	}
// 	idecl, ok := spec.Type.(*ast.InterfaceType)
// 	if !ok {
// 		return nil, fmt.Errorf("not an interface: %s", iface)
// 	}

// 	if idecl.Methods == nil {
// 		return nil, fmt.Errorf("empty interface: %s", iface)
// 	}

// 	var fns []Func
// 	for _, fndecl := range idecl.Methods.List {
// 		if len(fndecl.Names) == 0 {
// 			// Embedded interface: recurse
// 			embedded, err := funcs(p.fullType(fndecl.Type), srcDir)
// 			if err != nil {
// 				return nil, err
// 			}
// 			fns = append(fns, embedded...)
// 			continue
// 		}

// 		fn := p.funcsig(fndecl)
// 		fns = append(fns, fn)
// 	}
// 	return fns, nil
// }

const stub = "func ({{.Recv}}) {{.Name}}" +
	"({{range .Params}}{{.Name}} {{.Type}}, {{end}})" +
	"({{range .Res}}{{.Name}} {{.Type}}, {{end}})" +
	"{\n" + "panic(\"not implemented\")" + "}\n\n"

var tmpl = template.Must(template.New("test").Parse(stub))

// genStubs prints nicely formatted method stubs
// for fns using receiver expression recv.
// If recv is not a valid receiver expression,
// genStubs will panic.
func genStubs(recv string, fns []Func) []byte {
	var buf bytes.Buffer
	for _, fn := range fns {
		meth := Method{Recv: recv, Func: fn}
		tmpl.Execute(&buf, meth)
	}

	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	return pretty
}

// validReceiver reports whether recv is a valid receiver expression.
func validReceiver(recv string) bool {
	if recv == "" {
		// The parse will parse empty receivers, but we don't want to accept them,
		// since it won't generate a usable code snippet.
		return false
	}
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "", "package hack\nfunc ("+recv+") Foo()", 0)
	return err == nil
}

func parseFile(srcPath string, imp []byte) *ast.File {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcPath, imp, 0)
	if err != nil {
		log.Printf("ParseFile file error %s", err)
		os.Exit(2)
	}
	// log.Printf("f: %#v", pretty.Formatter(f))
	fmt.Printf("f: %+v", f)
	fmt.Printf("Decls: %# v", pretty.Formatter(f.Decls))
	fmt.Printf("Imports: %# v", pretty.Formatter(f.Imports))
	return f
}

func funcs(file *ast.File) {
	for _, decl := range file.Decls {
		switch typ := decl.(type) {
		case *ast.FuncDecl:
			fmt.Printf("func %s", typ.Name.Name)
			arguments(typ)

			fset := token.NewFileSet()

			var buf bytes.Buffer
			err := format.Node(&buf, fset, typ.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(buf.String())

			// fmt.Printf("%s", typ.Body)
			// fmt.Printf("params: %#v", arguments*(typ.Type.Params)
		}
	}
}

func arg(param *ast.Field) {
	id, typ := param.Names[0], param.Type
	fmt.Printf("id: %s -> ", id.Name)

	switch v := typ.(type) {
	case *ast.SelectorExpr:
		fmt.Printf("%s.%s", v.X, v.Sel)
	case *ast.StarExpr:
		switch s := v.X.(type) {
		case *ast.SelectorExpr:
			fmt.Printf("*%s.%s", s.X, s.Sel)
		default:
		}
		// fmt.Printf("*%s.%s", v.X, v.Sel)
	case *ast.Ellipsis:
		switch s := v.Elt.(type) {
		case *ast.SelectorExpr:
			fmt.Printf("...%s.%s", s.X, s.Sel)
		default:
		}
	default:
	}
	fmt.Printf("\n")
}

func arguments(fun *ast.FuncDecl) {
	params := fun.Type.Params
	for _, field := range params.List {
		arg(field)
		// log.Printf("%s", field.)
	}
}

type conditionChain struct {
}

type Condition interface {
	Cond(ast.Node) bool
}

type IsFunc struct {
}

func (s *isFunc) Cond(ast.Node) bool {
	x, ok := n.(*ast.FuncDecl)
	if !ok {
		return true
	}
	return false
}

func walkFunc(n ast.Node) {

}

func main() {
	flag.Parse()

	var code = []byte(`package hack
func (c *Clients) ListProject(ctx context.Context, in *types.User, opts ...grpc.CallOption) (*types1.Empty, error) {
	panic("some thing")
}
`)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "foo.go", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	rewriteFunc := func(n ast.Node) (ast.Node, bool) {
		x, ok := n.(*ast.FuncDecl)
		if !ok {
			return n, true
		}

		// change struct type name to "Bar"
		x.Name.Name = "Bar"
		return x, true
	}

	rewritten := astrewrite.Walk(file, rewriteFunc)

	var buf bytes.Buffer
	printer.Fprint(&buf, fset, rewritten)
	fmt.Println(buf.String())
	ast.Print(fset, file)
}

func fatal(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
