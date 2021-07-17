package apps

import (
	"fmt"
	"github.com/just1689/kubernetes-warm-images/client"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var namespacesFilename = "/config/list.spaces" //TODO: override in env var

func RunController() {
	logrus.Infoln(util.LogPrepend(1, "~~~ Started as Controller ~~~"))

	//Get namespaces to watch
	logrus.Infoln(util.LogPrepend(2, "getting namespaces to watch"))
	ns := getNamespacesToWatch()

	//Connect to the NATs server
	logrus.Infoln(util.LogPrepend(2, "connecting to pubSub"))
	pubSub := client.NewPubSubClient()

	//Connect to Kubernetes API
	logrus.Infoln(util.LogPrepend(2, "connecting to K8s"))
	k8sClient := client.NewK8sClient()

	//Subscribe to CREATE Pod
	images := k8sClient.WatchImages(ns)
	go func() {
		for image := range images {
			logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("Found new image:%s", image)))
			pubSub.Publish(image)
		}
	}()
	startHealthServer() //Blocking call
}

func getNamespacesToWatch() chan string {
	c, err := client.SplitFileBySpace(namespacesFilename)
	if err != nil {
		logrus.Panicln("could not get namespace list")
	}
	return c
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
