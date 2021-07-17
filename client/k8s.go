package client

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var emptyArr = []string{}

func NewK8sClient() *K8sClient {
	result := &K8sClient{}
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	result.clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorln("could not connect to Kubernetes API in-cluster. Exiting")
		logrus.Panicln(err.Error())
	}
	return result
}

type K8sClient struct {
	clientSet *kubernetes.Clientset
}

func (cl *K8sClient) WatchImages(namespaces chan string) chan string {
	result := make(chan string, 256)
	for namespace := range namespaces {
		if namespace == "*" {
			namespace = ""
		}
		logrus.Infoln("New watcher ns:", namespace)
		wi := cl.newWatcher(namespace)
		go handleWatch(namespace, wi, result)
	}
	return result
}

func handleWatch(namespace string, wi watch.Interface, imgChan chan string) {
	logrus.Infoln("handleWatch:", namespace)
	for event := range wi.ResultChan() {
		arrToChan(getImages(event), imgChan)
	}
}

func getImages(event watch.Event) []string {
	if event.Type != "ADDED" {
		return emptyArr
	}
	if p, ok := event.Object.(*v1.Pod); !ok {
		logrus.Errorln("could not cast as pod", event.Object)
		return emptyArr
	} else {
		result := make([]string, len(p.Spec.Containers))
		for i, container := range p.Spec.Containers {
			result[i] = container.Image
		}
		return result
	}
}

func (cl *K8sClient) newWatcher(namespace string) (wi watch.Interface) {
	logrus.Infoln("watching pods in ns:", namespace)
	var err error
	if wi, err = cl.clientSet.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{}); err != nil {
		logrus.Errorln(fmt.Sprintf("could not watch Pods in namespace `%s`", namespace))
		return
	}
	return
}

func arrToChan(arr []string, imgChan chan string) {
	for _, next := range arr {
		imgChan <- next
	}
}
