package ast

import (
	"maps"
	"slices"

	"go.ufukty.com/gommons/pkg/tree"
)

type Dimensional map[string]uint

// Gss
// TODO: update [safeEq] with Gss node prop value types
type (
	Display struct {
		Outside string
		Inside  string
	}
	Border struct {
		Color any // "inherit", "transparent", [color.RGBA]
		Style any // "inherit", gss.BorderStyle
		Width any // "none", [Dimensional]
	}
	BorderRadiuses struct {
		TopLeft, TopRight, BottomRight, BottomLeft any // "none", "inherit", [Dimensional]
	}
	Borders struct {
		Top, Right, Bottom, Left Border
	}
	Margin struct {
		Top, Right, Bottom, Left any // "inherit", [Dimensional]
	}
	Padding struct {
		Top, Right, Bottom, Left any // "inherit", [Dimensional]
	}
	Font struct {
		Family any // "inherit", []string
		Size   any // "inherit", [Dimensional]
		Weight any // "inherit", int
	}
	Text struct {
		Color         any // "inherit", "transparent", [color.RGBA]
		LineHeight    any // "inherit", [Dimensional]
		TextAlignment any // "inherit", [tokens.TextAlignment]
	}
	Dimensions struct {
		Height any // "auto", "min-content", "max-content", [Dimensional]
		Width  any // "auto", "min-content", "max-content", [Dimensional]
	}
	// TODO: handle shorthand syntaxes during parsing
	Styles struct {
		Dimensions      Dimensions
		Margin          Margin
		Padding         Padding
		Display         Display
		Text            Text
		Font            Font
		Border          Borders
		BorderRadiuses  BorderRadiuses
		BackgroundColor any // "inherit", "transparent", [color.RGBA]
	}
	Rule struct {
		Selector string
		Styles   *Styles
	}
	Gss struct {
		Rules []*Rule
	}
)

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
				ss = append(ss, tree.List(k, []string{v}))
			case interface{ String() string }:
				ss = append(ss, tree.List(k, []string{v.String()}))
			case interface{ Strings() []string }:
				ss = append(ss, tree.List(k, v.Strings()))
			default:
				ss = append(ss, tree.List(k, []string{"value of unknown type"}))
			}
		}
	}
	if len(ss) > 0 {
		return ss
	}
	return nil
}

//gomethods:exported
func (s Display) Strings() []string {
	return collect(map[string]any{
		"Outside": s.Outside,
		"Inside":  s.Inside,
	})
}

//gomethods:exported
func (s Border) Strings() []string {
	return collect(map[string]any{
		"Color":     s.Color,
		"Style":     s.Style,
		"Thickness": s.Width,
	})
}

//gomethods:exported
func (s BorderRadiuses) Strings() []string {
	return collect(map[string]any{
		"TopLeft":     s.TopLeft,
		"TopRight":    s.TopRight,
		"BottomRight": s.BottomRight,
		"BottomLeft":  s.BottomLeft,
	})
}

//gomethods:exported
func (s Borders) Strings() []string {
	return collect(map[string]any{
		"Top":    s.Top,
		"Right":  s.Right,
		"Bottom": s.Bottom,
		"Left":   s.Left,
	})
}

//gomethods:exported
func (s Margin) Strings() []string {
	return collect(map[string]any{
		"Top":    s.Top,
		"Right":  s.Right,
		"Bottom": s.Bottom,
		"Left":   s.Left,
	})
}

//gomethods:exported
func (s Padding) Strings() []string {
	return collect(map[string]any{
		"Top":    s.Top,
		"Right":  s.Right,
		"Bottom": s.Bottom,
		"Left":   s.Left,
	})
}

//gomethods:exported
func (s Font) Strings() []string {
	return collect(map[string]any{
		"Family": s.Family,
		"Size":   s.Size,
		"Weight": s.Weight,
	})
}

//gomethods:exported
func (s Text) Strings() []string {
	return collect(map[string]any{
		"Color":         s.Color,
		"LineHeight":    s.LineHeight,
		"TextAlignment": s.TextAlignment,
	})
}

//gomethods:exported // want `missing fields: Height, Width`
func (s Dimensions) Strings() []string {
	return nil
}

//gomethods:exported
func (s Styles) Strings() []string {
	return collect(map[string]any{
		"Dimensions":      s.Dimensions,
		"Margin":          s.Margin,
		"Padding":         s.Padding,
		"Display":         s.Display,
		"Text":            s.Text,
		"Font":            s.Font,
		"Border":          s.Border,
		"BorderRadiuses":  s.BorderRadiuses,
		"BackgroundColor": s.BackgroundColor,
	})
}

//gomethods:exported
func (r Rule) String() string {
	return tree.List(r.Selector, r.Styles.Strings())
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

//gomethods:all
func (s Display) IsEqual(y Display) bool {
	return s.Outside == y.Outside &&
		s.Inside == y.Inside
}

//gomethods:all
func (s Border) IsEqual(y Border) bool {
	return s.Color == y.Color &&
		s.Style == y.Style &&
		safeEq(s.Width, y.Width)
}

//gomethods:all
func (s BorderRadiuses) IsEqual(y BorderRadiuses) bool {
	return safeEq(s.TopLeft, y.TopLeft) &&
		safeEq(s.TopRight, y.TopRight) &&
		safeEq(s.BottomRight, y.BottomRight) &&
		safeEq(s.BottomLeft, y.BottomLeft)
}

//gomethods:all // want `missing fields: Bottom, Left, Right, Top`
func (s Borders) IsEqual(y Borders) bool {
	return false
}

//gomethods:all // want `missing field: Top`
func (s Margin) IsEqual(y Margin) bool {
	return safeEq(s.Right, y.Right) &&
		safeEq(s.Bottom, y.Bottom) &&
		safeEq(s.Left, y.Left)
}

//gomethods:all
func (s Padding) IsEqual(y Padding) bool {
	return safeEq(s.Top, y.Top) &&
		safeEq(s.Right, y.Right) &&
		safeEq(s.Bottom, y.Bottom) &&
		safeEq(s.Left, y.Left)
}

//gomethods:all
func (s Font) IsEqual(y Font) bool {
	return safeEq(s.Family, y.Family) &&
		safeEq(s.Size, y.Size) &&
		safeEq(s.Weight, y.Weight)
}

//gomethods:all
func (s Text) IsEqual(y Text) bool {
	return s.Color == y.Color &&
		safeEq(s.LineHeight, y.LineHeight) &&
		s.TextAlignment == y.TextAlignment
}

//gomethods:all
func (s Dimensions) IsEqual(y Dimensions) bool {
	return safeEq(s.Height, y.Height) &&
		safeEq(s.Width, y.Width)
}

//gomethods:all
func (s Styles) IsEqual(y Styles) bool {
	return s.Dimensions.IsEqual(y.Dimensions) &&
		s.Margin.IsEqual(y.Margin) &&
		s.Padding.IsEqual(y.Padding) &&
		s.Display.IsEqual(y.Display) &&
		s.Text.IsEqual(y.Text) &&
		s.Font.IsEqual(y.Font) &&
		s.Border.IsEqual(y.Border) &&
		s.BorderRadiuses.IsEqual(y.BorderRadiuses) &&
		s.BackgroundColor == y.BackgroundColor
}
