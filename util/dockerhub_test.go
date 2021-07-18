package util

import (
	"fmt"
	"testing"
)

func TestPrependImageWithRegistry(t *testing.T) {
	imageIn := "docker.io/library/nginx:latest"
	result := PrependImage(imageIn)
	expected := imageIn
	if expected != result {
		t.Error(fmt.Sprintf("expected '%s' and got '%s' instead", expected, result))
	}
}

func TestPrependImageWithoutRegistry(t *testing.T) {
	imageIn := "docker.io/library/nginx:latest"
	result := PrependImage(imageIn)
	expected := imageIn
	if expected != result {
		t.Error(fmt.Sprintf("expected '%s' and got '%s' instead", expected, result))
	}
}
