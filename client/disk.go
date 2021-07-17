package client

import (
	"fmt"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

func readFileToString(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Errorln(util.LogPrepend(3, err.Error()))
		return "", err
	}
	return string(b), nil
}

func SplitFileBySpace(file string) (res chan string, err error) {
	var b []byte
	b, err = ioutil.ReadFile(file)
	if err != nil {
		logrus.Errorln(util.LogPrepend(3, fmt.Sprintf("could not read file: %s", file)))
		logrus.Errorln(util.LogPrepend(3, err.Error()))
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
