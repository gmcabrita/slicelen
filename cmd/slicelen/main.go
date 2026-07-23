// Command slicelen runs the slice length indexing analyzer.
package main

import (
	slicelen "github.com/gmcabrita/go-slice-len-analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(slicelen.Analyzer)
}
