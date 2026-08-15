// Command slicelen runs the slice length indexing analyzer.
package main

import (
	"github.com/gmcabrita/slicelen"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(slicelen.Analyzer)
}
