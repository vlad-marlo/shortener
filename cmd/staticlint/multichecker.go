package main

import (
	"fmt"
	"regexp"

	"github.com/gostaticanalysis/nilerr"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/vlad-marlo/shortener/pkg/osexitanalyzer"
)

// main ...
func main() {
	var myChecks = []*analysis.Analyzer{
		osexitanalyzer.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		nilerr.Analyzer,
		bodyclose.Analyzer,
	}

	comp, err := regexp.Compile(`^SA[0-9]{4}$`)
	if err != nil {
		panic(fmt.Sprintf("regexp compile: %v", err))
	}

	for _, v := range staticcheck.Analyzers {
		if v.Analyzer.Name == "ST1020" || comp.MatchString(v.Analyzer.Name) {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	multichecker.Main(
		myChecks...,
	)
}
