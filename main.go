package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/go/ast/astutil"
)

const (
	defaultSourceIdent = "T__"
)

var (
	sourcePkg    = flag.String("src", "", "the package to generate code from")
	replace      = flag.String("replace", "", "a comma separated list of replacements to perform")
	excludeTypes = flag.String("exclude", "", "a comma separated list of types to exclude from generation, defaults to none excluded")
	includeTypes = flag.String("include", "", "a comma separated list of types to include in generation, defaults to all included")
)

func FilterIdents() func(c *astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		n := c.Node()
		switch n := n.(type) {
		case *ast.GenDecl:
			for _, spec := range n.Specs {
				if v, ok := spec.(*ast.TypeSpec); ok {
					if !isIncludedIdent(v.Name.Name) {
						c.Delete()
					}
				}
			}
		case *ast.FuncDecl:
			if n.Recv == nil {
				if !isIncludedIdent(n.Name.String()) {
					c.Delete()
				}
			} else {
				if !isIncludedIdent(fmt.Sprintf("%s.%s", n.Recv.List[0].Type, n.Name.String())) {
					c.Delete()
				}
			}
		}
		return true
	}
}

func ReplaceIdent(from, to string) func(c *astutil.Cursor) bool {
	return func(c *astutil.Cursor) bool {
		n := c.Node()
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

func isIncludedIdent(name string) bool {
	excludeIdents := strings.Split(*excludeTypes, ",")
	excludeIdents = append(excludeIdents, defaultSourceIdent)
	for _, ident := range excludeIdents {
		if ident == name {
			return false
		}
	}
	if *includeTypes == "" {
		return true
	}
	includeIdents := strings.Split(*includeTypes, ",")
	for _, ident := range includeIdents {
		if name == ident {
			return true
		}
	}
	return false
}

func findAndReplace(match, find, replace string) string {
	replaceLower := strcase.ToLowerCamel(replace)
	findLower := strcase.ToLowerCamel(find)
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
		for _, replacement := range getReplacements() {
			err := generateAll(wd, pkg.Name, replacement, *sourcePkg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type replacement struct{ from, to string }

func getReplacements() []replacement {
	replacementStrings := strings.Split(*replace, ",")
	var replacements []replacement
	for _, repStr := range replacementStrings {
		rep := strings.Split(repStr, "=")
		if len(rep) != 2 {
			continue
		}
		replacements = append(replacements, replacement{
			from: rep[0], to: rep[1],
		})
	}
	return replacements
}

func generateAll(path string, pkg string, typ replacement, declarations ...string) error {
	for i := range declarations {
		splitType := strings.Split(declarations[i], ":")
		var ipath, ident string
		if len(splitType) == 2 {
			ipath = splitType[0]
			ident = splitType[1]
		} else {
			ipath = declarations[i]
			ident = typ.from
		}
		err := generate(path, pkg, typ.to, ipath, ident)
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
	pkgs, err := parser.ParseDir(fset, resolvedDeclaration, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// Replace the package idents
	for _, p := range pkgs {
		// Replace all generic idents
		astutil.Apply(p, FilterIdents(), ReplaceIdent(ident, typ))
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
	err = format.Node(f, fset, file)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

func envOr(name, def string) string {
	s := os.Getenv(name)
	if s == "" {
		return def
	}
	return s
}

func resolveDeclarationPath(decl string) (string, error) {
	gopath := envOr("GOPATH", defaultGOPATH())
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
