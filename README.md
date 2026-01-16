# GoMethods

GoMethods is a linter plugin that warns when a method does not explicitly reference all required fields on its receiver struct. It is designed to warn developer when the field list of a struct changed without also updating its serializer, comparator etc. Opt-in per method via doc directives:

1. Require all receiver fields (exported + unexported)

   ```go
   //gomethods:all
   func (r Receiver) Validate() map[string]error
   ```

1. Require only exported receiver fields

   ```go
   //gomethods:exported
   func (r Receiver) String() string
   ```

Note that GoMethods only count **direct** selector usage rooted at the receiver. Passing the receiver to another function (e.g. `json.Marshal(r)`) does NOT count as referencing fields.

## Install

```sh
go install go.ufukty.com/gomethods@v0.1.0
which gomethods
```

## Example

Take a look at the [testfile](pkg/analyzer/testdata/gss.go):

```go
type (
  Borders struct {
    Top, Right, Bottom, Left Border
  }
  Margin struct {
    Top, Right, Bottom, Left any
  }
  Dimensions struct {
    Height     any
    Width      any
    unexported any
  }
)

//gomethods:exported
func (s Dimensions) Strings() []string {
  return nil
}

//gomethods:all
func (s Borders) IsEqual(y Borders) bool {
  return false
}

//gomethods:all
func (s Margin) IsEqual(y Margin) bool {
  return safeEq(s.Right, y.Right) &&
    safeEq(s.Bottom, y.Bottom) &&
    safeEq(s.Left, y.Left)
}
```

```sh
cd pkg/analyzer/testdata
gomethods .
pkg/analyzer/testdata/gss.go:173:1: missing fields: Height, Width
pkg/analyzer/testdata/gss.go:241:1: missing fields: Bottom, Left, Right, Top
pkg/analyzer/testdata/gss.go:246:1: missing field: Top
```
