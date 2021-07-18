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

func StrArrToCh(arr []string) chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		for _, next := range arr {
			if next != "" {
				result <- next
			}
		}
	}()
	return result
}
