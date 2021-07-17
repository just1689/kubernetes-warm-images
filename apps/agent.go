package apps

import (
	"fmt"
	"github.com/just1689/kubernetes-warm-images/client"
	"github.com/sirupsen/logrus"
	"os"
)

var HostIP = os.Getenv("HOST_IP")

func RunAgent() {
	logrus.Infoln(logPrepend(1, "~~~ Started as Agent ~~~"))
	listenForImages()
	logrus.Infoln(logPrepend(2, "starting health server"))
	startHealthServer()
}

func listenForImages() {
	logrus.Infoln(logPrepend(2, "new PubSub client"))
	ps := client.NewPubSubClient()
	logrus.Infoln(logPrepend(2, "start image puller"))
	pull := startDockerPuller()
	logrus.Infoln(logPrepend(2, "subscribing for images"))
	ps.Subscribe(func(image string) {
		logrus.Infoln(logPrepend(3, fmt.Sprintf("queue=%s", image)))
		pull <- image
	})
}

func startDockerPuller() chan string {
	in := make(chan string)
	go func() {
		for next := range in {
			client.PullDockerImageGo(next, logPrepend)
		}
	}()
	return in
}

func logPrepend(level int, text string) (result string) {
	result = ""
	if level == 1 {
		result += ">   "
	} else if level == 2 {
		result += ">>  "
	} else {
		result += "    "
	}
	result += fmt.Sprintf("%s %s: %s", HostIP, result, text)
	return
}
