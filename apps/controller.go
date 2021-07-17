package apps

import (
	"github.com/just1689/kubernetes-warm-images/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var namespacesFilename = "/config/list.spaces" //TODO: override in env var

func RunController() {
	logrus.Infoln("> Started as Controller")

	//Get namespaces to watch
	ns := getNamespacesToWatch()

	//Connect to the NATs server
	pubSub := client.NewPubSubClient()

	//Connect to Kubernetes API
	k8sClient := client.NewK8sClient()

	//Subscribe to CREATE Pod
	images := k8sClient.WatchImages(ns)
	go func() {
		for image := range images {
			logrus.Infoln("new image: ", image)
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
