package osexitanalyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

const basedir = "testdata/osexitanalyzer"

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, basedir+"/pkg1", Analyzer, "./pkg/osexitanalyzer/testdata/osexitanalyzer/pkg1")
}
