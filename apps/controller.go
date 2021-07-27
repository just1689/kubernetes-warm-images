package apps

import (
	"fmt"
	"github.com/just1689/kubernetes-warm-images/client"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var namespacesWatchFilename = util.StrOr(os.Getenv("WATCH_NAMESPACES_FILENAME"), "/config/list.spaces")
var namespacesIgnoreFilename = util.StrOr(os.Getenv("IGNORE_NAMESPACES_FILENAME"), "/config/ignore.spaces")

func RunController() {
	logrus.Infoln(util.LogPrepend(1, "~~~ Started as Controller ~~~"))

	pubSub := client.ConnectToNATs()

	go proxyImagesToQueue(pubSub)
	startHealthServer()
}

func proxyImagesToQueue(pubSub *client.PubSub) {
	//Get state
	//TODO: get state

	for image := range getImageStream() {
		logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("publishing image: %s", image)))
		pubSub.Publish(image)
	}
}

func getImageStream() chan string {
	//Get namespaces to watch & ignore
	logrus.Infoln(util.LogPrepend(2, "getting namespaces to watch & ignore"))
	ns := util.StrArrToCh(strings.Split(client.FileStrOrPanic(namespacesWatchFilename), " "))
	nsIgnore := strings.Split(client.FileStrOrPanic(namespacesIgnoreFilename), " ")
	//Connect to Kubernetes API
	logrus.Infoln(util.LogPrepend(2, "connecting to K8s"))
	k8s := client.ConnectToKubernetesAPI()
	//Subscribe to CREATE Pod
	return k8s.WatchEachNamespace(ns, nsIgnore)

}
