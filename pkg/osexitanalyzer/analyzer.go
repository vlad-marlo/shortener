package osexitanalyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for unchecked errors",
	Run: func(p *analysis.Pass) (interface{}, error) {
		for _, f := range p.Files {
			ast.Inspect(f, func(node ast.Node) bool {
				for _, c := range f.Comments {
					// if file was generated passes inspection
					if strings.Contains(c.Text(), "DO NOT EDIT") {
						return false
					}
				}
				switch x := node.(type) {
				case *ast.File:
					if x.Name.Name != "main" {
						return false
					}
				case *ast.SelectorExpr:
					if x.Sel.Name == "Exit" {
						p.Reportf(x.Pos(), "found call os.Exit() in main pkg")
					}
				}
				return true
			})
		}
		return nil, nil
	},
}