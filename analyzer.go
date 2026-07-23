// Package slicelen provides an analyzer that detects indexing a slice or array
// with its own length.
package slicelen

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

const diagnostic = "indexing a slice or array with len(...) always panics"

// Analyzer reports expressions of the form x[len(x)] when x is a slice or an
// array. It runs despite type errors because the compiler itself rejects
// array[len(array)] as a constant out-of-bounds index.
var Analyzer = &analysis.Analyzer{
	Name:             "slicelen",
	Doc:              "report indexing a slice or array with its own length",
	RunDespiteErrors: true,
	Run:              run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := inspector.New(pass.Files)
	inspect.Preorder([]ast.Node{(*ast.IndexExpr)(nil)}, func(node ast.Node) {
		index, ok := node.(*ast.IndexExpr)
		if !ok || !isSliceOrArray(pass.TypesInfo.TypeOf(index.X)) {
			return
		}

		call, ok := unparen(index.Index).(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return
		}

		function, ok := unparen(call.Fun).(*ast.Ident)
		if !ok || !isBuiltinLen(pass.TypesInfo, function) {
			return
		}

		if sameExpr(pass.TypesInfo, index.X, call.Args[0]) {
			pass.Report(analysis.Diagnostic{
				Pos:     index.Index.Pos(),
				End:     index.Index.End(),
				Message: diagnostic,
			})
		}
	})

	return nil, nil
}

func isSliceOrArray(typ types.Type) bool {
	if typ == nil {
		return false
	}

	switch typ.Underlying().(type) {
	case *types.Slice, *types.Array:
		return true
	default:
		return false
	}
}

func isBuiltinLen(info *types.Info, ident *ast.Ident) bool {
	builtin, ok := info.ObjectOf(ident).(*types.Builtin)
	return ok && builtin.Name() == "len"
}

func sameExpr(info *types.Info, left, right ast.Expr) bool {
	left = unparen(left)
	right = unparen(right)

	switch left := left.(type) {
	case *ast.Ident:
		right, ok := right.(*ast.Ident)
		if !ok {
			return false
		}

		leftObject := info.ObjectOf(left)
		return leftObject != nil && leftObject == info.ObjectOf(right)

	case *ast.SelectorExpr:
		right, ok := right.(*ast.SelectorExpr)
		if !ok || !sameExpr(info, left.X, right.X) {
			return false
		}

		leftObject := selectorObject(info, left)
		return leftObject != nil && leftObject == selectorObject(info, right)

	case *ast.IndexExpr:
		right, ok := right.(*ast.IndexExpr)
		return ok &&
			sameExpr(info, left.X, right.X) &&
			sameExpr(info, left.Index, right.Index)

	default:
		return false
	}
}

func selectorObject(info *types.Info, selector *ast.SelectorExpr) types.Object {
	if object := info.ObjectOf(selector.Sel); object != nil {
		return object
	}
	if selection := info.Selections[selector]; selection != nil {
		return selection.Obj()
	}
	return nil
}

func unparen(expr ast.Expr) ast.Expr {
	for {
		paren, ok := expr.(*ast.ParenExpr)
		if !ok {
			return expr
		}
		expr = paren.X
	}
}
