package client

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func ReadFileToString(file string) (result string, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(file); err != nil {
		logrus.Panicln(fmt.Sprintf("could not read file: %s with error: %s", file, err))
		return
	}
	result = string(b)
	return
}

func fileStrOrPanic(path string) (str string) {
	var err error
	if str, err = ReadFileToString(path); err != nil || str == "" {
		logrus.Panicln(fmt.Sprintf("no ne string found '%s' found at %s", str, path))
	}
	return
}
