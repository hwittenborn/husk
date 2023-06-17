## Coding style
### `Debug` formatting
In structs that contain Go pointers (a private `ptr` or `dptr` field), `Debug` should be manually implemented, with the pointer field ommited.

Some structs also have other fields added (such as with `syntax::Pos`) - do whatever feels right and we can go from there.

### `ptr` vs `dptr` naming
In structs containing Go pointers, they'll have the pointer registered under a struct field named either `ptr` or `dptr`:

- `ptr` should be used when referencing a Go pointer directly (i.e. a pointer to `ParseError`)
- `dptr` should be used when referencing a *reference* of a Go pointer (i.e. a pointer of `*Pointer`)

This rule has been put in place for clarity when running FFI code between Go and Rust
