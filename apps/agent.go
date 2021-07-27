package apps

import (
	"fmt"
	"github.com/just1689/kubernetes-warm-images/client"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
)

var ImageQueueCacheSize = 256 // TODO: consider making env var & Helm value

func RunAgent() {
	logrus.Infoln(util.LogPrepend(1, "~~~ Started as Agent ~~~"))
	listenForImages()
	startHealthServer()
}

func listenForImages() {
	ps := client.ConnectToNATs()
	ps.Subscribe(util.FuncForEachStr(imagePullCh(), func(next string) {
		logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("enqueuing image: '%s'", next)))
	}))
}

func imagePullCh() (result chan string) {
	logrus.Infoln(util.LogPrepend(2, "start image puller"))
	result = make(chan string, ImageQueueCacheSize)
	go util.FuncForEach(result, client.ConnectToDocker())
	return
}
