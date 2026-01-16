package ast

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
		Height     any // "auto", "min-content", "max-content", [Dimensional]
		Width      any // "auto", "min-content", "max-content", [Dimensional]
		unexported any
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
	return treeList(r.Selector, r.Styles.Strings())
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

//gomethods:exported
func (s Dimensions) IsEqual(y Dimensions) bool {
	return safeEq(s.Height, y.Height) &&
		safeEq(s.Width, y.Width)
}
