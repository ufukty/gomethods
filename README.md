# gomethods

It is a linter-like program that assures all annotated methods mention their receivers' all fields. It is designed to warn developer when the field list of a struct changed without also updating its serializer, comparator etc.

```go
type Multiplication struct {
  Left, Right int
}

//gomethods:all-fields
func (m Multiplication) String() string {
  return fmt.Sprintf("%d + %d", m.Left, m.Right)
}
```

```sh
gomethods run
echo $?
0
```