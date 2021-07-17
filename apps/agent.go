package apps

import (
	"fmt"
	"github.com/just1689/kubernetes-warm-images/client"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
)

func RunAgent() {
	logrus.Infoln(util.LogPrepend(1, "~~~ Started as Agent ~~~"))
	listenForImages()
	logrus.Infoln(util.LogPrepend(2, "starting health server"))
	startHealthServer()
}

func listenForImages() {
	logrus.Infoln(util.LogPrepend(2, "new PubSub client"))
	ps := client.NewPubSubClient()
	logrus.Infoln(util.LogPrepend(2, "start image puller"))
	pull := startDockerPuller()
	logrus.Infoln(util.LogPrepend(2, "subscribing for images"))
	ps.Subscribe(func(image string) {
		logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("queue=%s", image)))
		pull <- image
	})
}

func startDockerPuller() chan string {
	in := make(chan string)
	go func() {
		for next := range in {
			client.PullDockerImageGo(next, util.LogPrepend)
		}
	}()
	return in
}
