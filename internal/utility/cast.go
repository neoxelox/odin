package utility

import (
	"fmt"
	"strings"
)

const SIZE_BASE = 1024

func SizeToString(size int) string {
	if size < SIZE_BASE {
		return fmt.Sprintf("%dB", size)
	}

	div, exp := int64(SIZE_BASE), 0
	for n := size / SIZE_BASE; n >= SIZE_BASE; n /= SIZE_BASE {
		div *= SIZE_BASE
		exp++
	}

	number := fmt.Sprintf("%.1f", float64(size)/float64(div))
	exponent := "KMGTPE"[exp]

	return fmt.Sprintf("%s%cB", strings.TrimRight(strings.TrimRight(number, "0"), "."), exponent)
}
