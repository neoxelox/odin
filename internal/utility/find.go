package utility

func StringIn(element string, slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if element == slice[i] {
			return true
		}
	}

	return false
}
