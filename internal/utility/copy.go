package utility

import (
	"time"

	"github.com/aodin/date"
)

func CopyInt(src *int) *int {
	if src == nil {
		return nil
	}

	dst := *src

	return &dst
}

func CopyIntSlice(src *[]int) *[]int {
	if src == nil {
		return nil
	}

	dst := make([]int, len(*src))
	copy(dst, *src)

	return &dst
}

func CopyIntMap(src *map[string]int) *map[string]int {
	if src == nil {
		return nil
	}

	dst := make(map[string]int, len(*src))
	for k, v := range *src {
		dst[k] = v
	}

	return &dst
}

func CopyIntSliceMap(src *map[string][]int) *map[string][]int {
	if src == nil {
		return nil
	}

	dst := make(map[string][]int, len(*src))
	for k, v := range *src {
		dst[k] = *CopyIntSlice(&v)
	}

	return &dst
}

func CopyString(src *string) *string {
	if src == nil {
		return nil
	}

	dst := *src

	return &dst
}

func CopyStringSlice(src *[]string) *[]string {
	if src == nil {
		return nil
	}

	dst := make([]string, len(*src))
	copy(dst, *src)

	return &dst
}

func CopyStringMap(src *map[string]string) *map[string]string {
	if src == nil {
		return nil
	}

	dst := make(map[string]string, len(*src))
	for k, v := range *src {
		dst[k] = v
	}

	return &dst
}

func CopyStringSliceMap(src *map[string][]string) *map[string][]string {
	if src == nil {
		return nil
	}

	dst := make(map[string][]string, len(*src))
	for k, v := range *src {
		dst[k] = *CopyStringSlice(&v)
	}

	return &dst
}

func CopyBool(src *bool) *bool {
	if src == nil {
		return nil
	}

	dst := *src

	return &dst
}

func CopyBoolSlice(src *[]bool) *[]bool {
	if src == nil {
		return nil
	}

	dst := make([]bool, len(*src))
	copy(dst, *src)

	return &dst
}

func CopyBoolMap(src *map[string]bool) *map[string]bool {
	if src == nil {
		return nil
	}

	dst := make(map[string]bool, len(*src))
	for k, v := range *src {
		dst[k] = v
	}

	return &dst
}

func CopyBoolSliceMap(src *map[string][]bool) *map[string][]bool {
	if src == nil {
		return nil
	}

	dst := make(map[string][]bool, len(*src))
	for k, v := range *src {
		dst[k] = *CopyBoolSlice(&v)
	}

	return &dst
}

func CopyTime(src *time.Time) *time.Time {
	if src == nil {
		return nil
	}

	dstYear, dstMonth, dstDay := src.Date()
	dstHour, dstMin, dstSec := src.Clock()
	dstNsec := src.Nanosecond()
	dstLocation := *src.Location()

	dst := time.Date(dstYear, dstMonth, dstDay, dstHour, dstMin, dstSec, dstNsec, &dstLocation)

	return &dst
}

func CopyTimeSlice(src *[]time.Time) *[]time.Time {
	if src == nil {
		return nil
	}

	dst := make([]time.Time, len(*src))
	for i := 0; i < len(dst); i++ {
		dst[i] = *CopyTime(&(*src)[i])
	}

	return &dst
}

func CopyTimeMap(src *map[string]time.Time) *map[string]time.Time {
	if src == nil {
		return nil
	}

	dst := make(map[string]time.Time, len(*src))
	for k, v := range *src {
		dst[k] = *CopyTime(&v)
	}

	return &dst
}

func CopyTimeSliceMap(src *map[string][]time.Time) *map[string][]time.Time {
	if src == nil {
		return nil
	}

	dst := make(map[string][]time.Time, len(*src))
	for k, v := range *src {
		dst[k] = *CopyTimeSlice(&v)
	}

	return &dst
}

func CopyDate(src *date.Date) *date.Date {
	if src == nil {
		return nil
	}

	dst := date.FromTime(src.Time)

	return &dst
}

func CopyDateSlice(src *[]date.Date) *[]date.Date {
	if src == nil {
		return nil
	}

	dst := make([]date.Date, len(*src))
	for i := 0; i < len(dst); i++ {
		dst[i] = *CopyDate(&(*src)[i])
	}

	return &dst
}

func CopyDateMap(src *map[string]date.Date) *map[string]date.Date {
	if src == nil {
		return nil
	}

	dst := make(map[string]date.Date, len(*src))
	for k, v := range *src {
		dst[k] = *CopyDate(&v)
	}

	return &dst
}

func CopyDateSliceMap(src *map[string][]date.Date) *map[string][]date.Date {
	if src == nil {
		return nil
	}

	dst := make(map[string][]date.Date, len(*src))
	for k, v := range *src {
		dst[k] = *CopyDateSlice(&v)
	}

	return &dst
}
