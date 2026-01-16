package gomethods

import (
	"go/ast"
	"go/types"
	"iter"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// function => signature => receiver => type => indirect => underlying
func receiverStruct(pass *analysis.Pass, fd *ast.FuncDecl) *types.Struct {
	if o, ok := pass.TypesInfo.Defs[fd.Name]; ok {
		if f, ok := o.(*types.Func); ok {
			if s, ok := f.Type().(*types.Signature); ok {
				if r := s.Recv(); r != nil {
					t := r.Type()
					if p, ok := t.(*types.Pointer); ok {
						t = p.Elem()
					}
					t = t.Underlying()
					if s, ok := t.(*types.Struct); ok {
						return s
					}
				}
			}
		}
	}
	return nil
}

var directivePattern = regexp.MustCompile(`^\s*//\s*gomethods:\s*(all|exported)\s*(?:// want .*)?$`)

func parseMode(cg *ast.CommentGroup) (string, analysis.Range, bool) {
	for _, c := range cg.List {
		if ms := directivePattern.FindStringSubmatch(c.Text); len(ms) > 1 && ms[1] != "" {
			return ms[1], c, true
		}
	}
	return "", nil, false
}

func fields(s *types.Struct) iter.Seq[*types.Var] {
	return func(yield func(*types.Var) bool) {
		for i := 0; i < s.NumFields(); i++ {
			if !yield(s.Field(i)) {
				return
			}
		}
	}
}

func filterFields(s *types.Struct, mode string) []string {
	var fs []string
	for f := range fields(s) {
		if f.Name() != "_" && (mode == "all" || f.Exported()) {
			fs = append(fs, f.Name())
		}
	}
	return fs
}

func isReceiverExpr(e ast.Expr, recvName string) bool {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name == recvName
	case *ast.ParenExpr:
		return isReceiverExpr(x.X, recvName)
	case *ast.StarExpr:
		return isReceiverExpr(x.X, recvName) // (*r).Field
	default:
		return false
	}
}

func receiverName(fd *ast.FuncDecl) string {
	if fd.Recv != nil && fd.Recv.List != nil && len(fd.Recv.List) > 0 {
		recv := fd.Recv.List[0]
		if recv.Names != nil && len(recv.Names) > 0 && recv.Names[0] != nil {
			return recv.Names[0].Name
		}
	}
	return ""
}

func listMentionedFields(fd *ast.FuncDecl) map[string]any {
	recvName := receiverName(fd)
	mentioned := map[string]any{}
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if ok && sel.Sel != nil && isReceiverExpr(sel.X, recvName) {
			mentioned[sel.Sel.Name] = nil
		}
		return true
	})
	return mentioned
}

func findMissing(required []string, mentioned map[string]any) []string {
	var missing []string
	for _, f := range required {
		if _, ok := mentioned[f]; !ok {
			missing = append(missing, f)
		}
	}
	return missing
}

func report(pass *analysis.Pass, pos analysis.Range, missing []string) {
	slices.Sort(missing)
	pattern := "missing field: %s"
	if len(missing) > 1 {
		pattern = "missing fields: %s"
	}
	pass.ReportRangef(pos, pattern, strings.Join(missing, ", "))
}

func inspectMethod(pass *analysis.Pass, fd *ast.FuncDecl) {
	if recvStruct := receiverStruct(pass, fd); recvStruct != nil {
		m, pos, _ := parseMode(fd.Doc)
		required := filterFields(recvStruct, m)
		if len(required) > 0 {
			missing := findMissing(required, listMentionedFields(fd))
			if len(missing) > 0 {
				report(pass, pos, missing)
			}
		}
	}
}

func listFuncDecls(fs []*ast.File) []*ast.FuncDecl {
	fds := []*ast.FuncDecl{}
	for _, f := range fs {
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok {
				fds = append(fds, fd)
			}
		}
	}
	return fds
}

// lists the methods with non-empty bodies and `//gomethods:` doc-comment
func listMethods(fs []*ast.File) []*ast.FuncDecl {
	fds := listFuncDecls(fs)
	ffds := make([]*ast.FuncDecl, 0, len(fds))
	for _, fd := range fds {
		if fd.Doc != nil {
			_, _, ok := parseMode(fd.Doc)
			if ok && fd.Recv != nil && len(fd.Recv.List) == 1 && fd.Body != nil {
				ffds = append(ffds, fd)
			}
		}
	}
	return ffds
}

func Inspect(pass *analysis.Pass) (any, error) {
	for _, fd := range listMethods(pass.Files) {
		inspectMethod(pass, fd)
	}
	return nil, nil
}
