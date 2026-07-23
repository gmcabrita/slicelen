# go-slice-len-analyzer

A focused [`go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis) analyzer that reports indexing a slice or array with its own length:

```go
value := s[len(s)] // indexing a slice or array with len(...) always panics
```

Valid indices end at `len(s)-1`, so this expression always panics at runtime. The analyzer compares resolved Go objects, handles parenthesized, selector, and nested index expressions, and verifies that `len` is the builtin function.

It intentionally does not perform general bounds or fencepost analysis. Expressions such as `s[len(s)+1]`, `s[len(s)-1]`, and `s[i+1]` are outside its scope.

## Run in another project

Run it directly from the target project's root:

```sh
go run github.com/gmcabrita/go-slice-len-analyzer/cmd/slicelen@latest ./...
```

Alternatively, install the command and reuse it across projects:

```sh
go install github.com/gmcabrita/go-slice-len-analyzer/cmd/slicelen@latest
cd /path/to/project
slicelen ./...
```

Arguments are standard Go package patterns. Use `./...` for the current module or specify individual packages.

## Embed the analyzer

The exported `slicelen.Analyzer` can be included in a custom multi-analyzer driver:

```go
package main

import (
	"github.com/gmcabrita/go-slice-len-analyzer"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(slicelen.Analyzer)
}
```

## Development

```sh
mise run check
```
