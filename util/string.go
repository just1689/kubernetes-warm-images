package util

func StrOr(first, second string) string {
	if first != "" {
		return first
	}
	return second
}

func StrExistsIn(needle string, haystack []string) bool {
	if len(haystack) == 0 {
		return false
	}
	for _, next := range haystack {
		if needle == next {
			return true
		}
	}
	return false
}
