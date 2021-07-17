package client

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

func readFileToString(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Errorln(err)
		return "", err
	}
	return string(b), nil
}

func SplitFileBySpace(file string) (res chan string, err error) {
	var b []byte
	b, err = ioutil.ReadFile(file)
	if err != nil {
		logrus.Errorln("could not read file ", file)
		logrus.Errorln(err)
		return
	}
	res = make(chan string)
	go func() {
		for _, next := range strings.Split(string(b), " ") {
			res <- next
		}
		close(res)
	}()
	return
}
