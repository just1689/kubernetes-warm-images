package client

import (
	"context"
	"fmt"
	"github.com/just1689/kubernetes-warm-images/util"
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
		logrus.Errorln(util.LogPrepend(3, "could not connect to Kubernetes API in-cluster. Exiting"))
		logrus.Panicln(err.Error())
	}
	return result
}

type K8sClient struct {
	clientSet *kubernetes.Clientset
}

func (cl *K8sClient) WatchImages(namespaces chan string, nsIgnore []string) chan string {
	result := make(chan string, 256)
	for namespace := range namespaces {
		if namespace == "*" {
			namespace = ""
		}
		logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("New watcher ns: '%s'", namespace)))
		wi := cl.newWatcher(namespace)
		go handleWatch(namespace, wi, result, nsIgnore)
	}
	return result
}

func handleWatch(namespace string, wi watch.Interface, imgChan chan string, nsIgnore []string) {
	logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("handleWatch namespace: '%s'", namespace)))
	for event := range wi.ResultChan() {
		util.StrArrToChan(getImages(event, nsIgnore), imgChan)
	}
}

func getImages(event watch.Event, nsIgnore []string) []string {
	if event.Type != "ADDED" {
		return emptyArr
	}
	if p, ok := event.Object.(*v1.Pod); !ok {
		logrus.Errorln(util.LogPrepend(3, fmt.Sprintf("could not cast as pod %s", event.Object)))
		return emptyArr
	} else {
		if util.StrExistsIn(p.ObjectMeta.Namespace, nsIgnore) {
			logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("ignoring pod: '%s' as namespace: '%s'", p.ObjectMeta.Name, p.ObjectMeta.Namespace)))
			return emptyArr
		}
		result := make([]string, len(p.Spec.Containers))
		for i, container := range p.Spec.Containers {
			result[i] = container.Image
		}
		return result
	}
}

func (cl *K8sClient) newWatcher(namespace string) (wi watch.Interface) {
	logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("watching pods in ns: '%s'", namespace)))
	var err error
	if wi, err = cl.clientSet.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{}); err != nil {
		logrus.Errorln(util.LogPrepend(3, fmt.Sprintf("could not watch Pods in namespace `%s`", namespace)))
		return
	}
	return
}
