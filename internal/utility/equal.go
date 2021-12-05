package utility

import (
	"github.com/scylladb/go-set/strset"
)

func EqualStringSlice(first *[]string, second *[]string) bool {
	if first == nil && second == nil {
		return true
	}

	if !(first != nil && second != nil) {
		return false
	}

	if len(*first) != len(*second) {
		return false
	}

	setFirst := strset.New((*first)...)
	setSecond := strset.New((*second)...)

	return strset.SymmetricDifference(setFirst, setSecond).IsEmpty()
}
