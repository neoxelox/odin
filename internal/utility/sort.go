package utility

import (
	"reflect"
)

// TODO: Discuss whether to make it a safe function
func EqualSort(primary interface{}, secondary interface{}, equal func(i int, j int) bool) {
	primarySlice := reflect.ValueOf(primary)
	lenSlice := primarySlice.Len()
	swap := reflect.Swapper(secondary)

	var i, j int
	for i = 0; i < lenSlice; i++ {
		for j = i; j < lenSlice && !equal(i, j); j++ {}
		swap(i, j)
	}
}
