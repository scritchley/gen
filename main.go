package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path"
	// "path/filepath"
	"strings"
)

const (
	defaultSourceIdent = "T__"
)

var (
	sourcePkg        = flag.String("src", "", "the package to generate code from")
	sourceIdent      = flag.String("ident", defaultSourceIdent, "the source ident to use for replacement")
	destinationIdent = flag.String("dest", "", "the destination ident to use for replacement")
)

func ReplaceIdent(from, to string) func(n ast.Node) bool {
	return func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			n.Name = findAndReplace(n.Name, from, to)
		case *ast.CommentGroup:
			for i := range n.List {
				n.List[i].Text = findAndReplace(n.List[i].Text, from, to)
			}
		case *ast.Comment:
			n.Text = findAndReplace(n.Text, from, to)
		}
		return true
	}
}

func excludeStub(f os.FileInfo) bool {
	return f.Name() != "stub.go" && f.Name() != "stub_test.go"
}

func findAndReplace(match, find, replace string) string {
	replaceLower := strings.ToLower(replace)
	findLower := strings.ToLower(find)
	if match == find {
		return replace
	}
	if strings.Contains(match, find) {
		return strings.Replace(match, find, replace, -1)
	}
	if strings.Contains(match, findLower) {
		return strings.Replace(match, findLower, replaceLower, -1)
	}
	return match
}

func match(s string) bool {
	return strings.Contains(s, "@")
}

func trim(s string) string {
	return strings.TrimPrefix(s, "// @")
}

func run() error {

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, wd, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		ast.Inspect(pkg, func(n ast.Node) bool {
			switch t := n.(type) {
			case *ast.TypeSpec:
				if t.Name.Name == *destinationIdent {
					generateAll(wd, pkg.Name, t.Name.Name, *sourcePkg)
				}
			}
			return true
		})
	}
	return nil

}

func generateAll(path string, pkg, typ string, declarations ...string) error {
	for i := range declarations {
		splitType := strings.Split(declarations[i], ":")
		var ipath, ident string
		if len(splitType) == 2 {
			ipath = splitType[0]
			ident = splitType[1]
		} else {
			ipath = declarations[i]
			ident = *sourceIdent
		}
		err := generate(path, pkg, typ, ipath, ident)
		if err != nil {
			return err
		}
	}
	return nil
}

func generate(path string, pkg, typ, declaration, ident string) error {
	log.Printf("generating %s: %s -> %s\n", declaration, ident, typ)
	resolvedDeclaration, err := resolveDeclarationPath(declaration)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, resolvedDeclaration, excludeStub, parser.ParseComments)
	if err != nil {
		return err
	}
	// Replace the package idents
	for _, p := range pkgs {
		// Replace all generic idents
		ast.Inspect(p, ReplaceIdent(ident, typ))
		// Replace the package name
		p.Name = pkg
		// Print the generated code
		err := writeAllFiles(path, fset, pkg, typ, p.Files)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeAllFiles(path string, fset *token.FileSet, pkg, typ string, files map[string]*ast.File) error {
	for filename := range files {
		files[filename].Name.Name = pkg
		err := writeFile(path, fset, typ, filename, files[filename])
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(outputPath string, fset *token.FileSet, typ, filename string, file *ast.File) error {
	filename = fmt.Sprintf("%s_%s", strings.ToLower(typ), path.Base(filename))
	fn := path.Join(outputPath, filename)
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	err = printer.Fprint(f, fset, file)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func resolveDeclarationPath(decl string) (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", fmt.Errorf("GOPATH not set")
	}
	return path.Join(gopath, "src", decl), nil
}

func main() {
	flag.Parse()
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
