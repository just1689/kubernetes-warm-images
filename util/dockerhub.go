package util

import "strings"

func PrependImage(image string) string {
	if !strings.Contains(image, ".") {
		image = "docker.io/library/" + image
	}
	return image
}
