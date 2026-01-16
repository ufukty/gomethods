package analyzer

import (
	"go.ufukty.com/gomethods/internal/gomethods"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "gomethods",
		Doc:      "method linter to remind missing fields",
		Run:      gomethods.Inspect,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}
