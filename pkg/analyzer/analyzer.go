package analyzer

import (
	"go.ufukty.com/golistics/internal/golistics"
	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "golistics",
		Doc:  "method linter to remind missing fields",
		Run:  golistics.Inspect,
	}
}
