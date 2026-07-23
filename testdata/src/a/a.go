package a

type collection struct {
	items []int
	other []int
}

type namedSlice []int

type namedArray [2]int

func findings(s []int, arr [2]int, obj collection, matrix [][]int, i int, ns namedSlice, na namedArray) {
	_ = s[len(s)]                 // want "indexing a slice or array with len\\(\\.\\.\\.\\) always panics"
	_ = arr[len(arr)]             // want "indexing a slice or array with len\\(\\.\\.\\.\\) always panics"
	_ = obj.items[len(obj.items)] // want "indexing a slice or array with len\\(\\.\\.\\.\\) always panics"
	_ = matrix[i][len(matrix[i])] // want "indexing a slice or array with len\\(\\.\\.\\.\\) always panics"
	_ = s[len(s)]                 // want "indexing a slice or array with len\\(\\.\\.\\.\\) always panics"
	_ = ns[len(ns)]               // want "indexing a slice or array with len\\(\\.\\.\\.\\) always panics"
	_ = na[len(na)]               // want "indexing a slice or array with len\\(\\.\\.\\.\\) always panics"
}

func shadowedBuiltin(s []int) {
	len := func([]int) int { return 0 }
	_ = s[len(s)]
}

func noFindings(s, other []int, obj collection, i int) {
	_ = s[len(other)]
	_ = obj.items[len(obj.other)]
	_ = s[len(s):]
	_ = s[:]
	_ = s[i]
}

func unsupported(s []int, m map[int]int, text string) {
	_ = s[len(s)+1]
	_ = s[len(s)-1]
	_ = m[len(m)]
	_ = text[len(text)]
}
