package analyzer

import (
	"go.ufukty.com/golistics/internal/golistics"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "golistics",
		Doc:      "method linter to remind missing fields",
		Run:      golistics.Inspect,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}
