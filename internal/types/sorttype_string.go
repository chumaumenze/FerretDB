// Code generated by "stringer -linecomment -type SortType"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Ascending-1]
	_ = x[Descending - -1]
}

const (
	_SortType_name_0 = "Descending"
	_SortType_name_1 = "Ascending"
)

func (i SortType) String() string {
	switch {
	case i == -1:
		return _SortType_name_0
	case i == 1:
		return _SortType_name_1
	default:
		return "SortType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}