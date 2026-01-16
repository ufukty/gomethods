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
