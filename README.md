[![GoDoc](https://godoc.org/github.com/KEINOS/go-errors?status.svg)](http://godoc.org/github.com/KEINOS/go-errors)
[![Report card](https://goreportcard.com/badge/github.com/KEINOS/go-errors)](https://goreportcard.com/report/github.com/KEINOS/go-errors)

# go-errors

"[github.com/KEINOS/go-errors](https://github.com/KEINOS/go-errors)" is a fork/replacement of [Dave Cheney](https://github.com/davecheney)'s unfortunately deprecated "[github.com/pkg/errors](https://github.com/pkg/errors)" awesome package. It provides simple error handling primitives.

## Usage

```go
// Download package
go get github.com/KEINOS/go-errors
```

```go
// Import module
import "github.com/KEINOS/go-errors/errors"
```

```diff
// Migration from "github.com/pkg/errors"
-import "github.com/pkg/errors"
+import "github.com/KEINOS/go-errors/errors"
```

## Note

If your main purpose is to wrap errors using `errors.Wrap()` or `errors.Wrapf()` and you do not need a stack trace, use the standard `error` package as below.

```go
func Wrap(err *error, format string, args ...any) {
    if *err != nil {
        *err = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *err)
    }
}
```

## Roadmap

There is no significant roadmap at this time. Other than updates to related modules and support for reporting with reproducible code, there are no plans to add new feature.

## Contributing

- Brach to PR: `main`
- As long as it keeps the backward compatibility, and the tests pass with 100% coverage, any PR for the better is welcome.
- Help wanted:
  - Refactoring
  - More examples
  - More documentation or better documentation

## License & Authors

- BSD-2-Clause
- Copyright (c) Dave Cheney, KEINOS and the contributors
