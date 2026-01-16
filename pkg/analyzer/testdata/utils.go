package ast

import (
	"maps"
	"slices"
)

// local stand-in for go.ufukty.com/gommons/pkg/tree.List
func treeList(_ ...any) string { return "" }

func isZero[T comparable](t T) bool {
	var z T
	return t == z
}

func collect(s map[string]any) []string {
	ss := []string{}
	for _, k := range slices.Sorted(maps.Keys(s)) {
		if !isZero(s[k]) {
			switch v := s[k].(type) {
			case string:
				ss = append(ss, treeList(k, []string{v}))
			case interface{ String() string }:
				ss = append(ss, treeList(k, []string{v.String()}))
			case interface{ Strings() []string }:
				ss = append(ss, treeList(k, v.Strings()))
			default:
				ss = append(ss, treeList(k, []string{"value of unknown type"}))
			}
		}
	}
	if len(ss) > 0 {
		return ss
	}
	return nil
}

// Use the first output for equality when second is true.
func try[T any](x, y any, eq func(x, y T) bool) (bool, bool) {
	xT, xOk := x.(T)
	yT, yOk := y.(T)
	if xOk && yOk {
		return eq(xT, yT), true
	}
	return false, xOk || yOk
}

// TODO: Try every possible incompatible operand type before
// using the built-in equality operator.
func safeEq(a, b any) bool {
	if r, ok := try(a, b, maps.Equal[Dimensional]); ok {
		return r
	}
	if r, ok := try(a, b, slices.Equal[[]string]); ok {
		return r
	}
	return a == b
}
