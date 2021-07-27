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

func StrArrToChan(arr []string, c chan string) {
	if len(arr) == 0 {
		return
	}
	for _, next := range arr {
		if next == "" {
			continue
		}
		c <- next
	}
}

func FuncForEach(in chan string, f func(string)) {
	for next := range in {
		f(next)
	}
}

func FuncForEachStr(in chan string, logger func(in string)) func(next string) {
	return func(next string) {
		logger(next)
		in <- next
	}
}
