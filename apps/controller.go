package apps

import (
	"fmt"
	"github.com/just1689/kubernetes-warm-images/client"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var namespacesWatchFilename = util.StrOr(os.Getenv("WATCH_NAMESPACES_FILENAME"), "/config/list.spaces")
var namespacesIgnoreFilename = util.StrOr(os.Getenv("IGNORE_NAMESPACES_FILENAME"), "/config/ignore.spaces")

func RunController() {
	logrus.Infoln(util.LogPrepend(1, "~~~ Started as Controller ~~~"))

	//Connect to the NATs server
	logrus.Infoln(util.LogPrepend(2, "connecting to pubSub"))
	pubSub := client.NewPubSubClient()

	//Subscribe to CREATE Pod
	go imagesToQueue(pubSub)
	startHealthServer() //Blocking call
}

func imagesToQueue(pubSub *client.PubSub) {
	for image := range getImageStream() {
		logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("Found new image:%s", image)))
		pubSub.Publish(image)
	}
}

func getImageStream() chan string {
	//Get namespaces to watch & ignore
	logrus.Infoln(util.LogPrepend(2, "getting namespaces to watch & ignore"))
	ns, nsIgnore := getNamespacesToWatch(), getNamespacesToIgnore()
	//Connect to Kubernetes API
	logrus.Infoln(util.LogPrepend(2, "connecting to K8s"))
	k8sClient := client.NewK8sClient()
	//Subscribe to CREATE Pod
	return k8sClient.WatchImages(ns, nsIgnore)

}

func getNamespacesToWatch() chan string {
	var s string
	var err error
	s, err = client.ReadFileToString(namespacesWatchFilename)
	if err != nil {
		logrus.Panicln("could not get namespace list to watch")
	}
	return util.StrArrToCh(strings.Split(s, " "))
}

func getNamespacesToIgnore() []string {
	var s string
	var err error
	s, err = client.ReadFileToString(namespacesIgnoreFilename)
	if err != nil {
		logrus.Panicln("could not get namespace list to ignore")
	}
	return strings.Split(s, " ")
}

func startHealthServer() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		logrus.Panicln("could not find env var LISTEN_ADDR. Exiting")
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//TODO: check the health is OK
		if false {
			errMsg := "" // TODO: error message
			http.Error(writer, errMsg, http.StatusInternalServerError)
			return
		}
		writer.Write([]byte("{\"ok\":true}"))
	})
	http.ListenAndServe(listenAddr, nil)
}
