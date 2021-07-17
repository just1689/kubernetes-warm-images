package util

func StrOr(first, second string) string {
	if first != "" {
		return first
	}
	return second
}
