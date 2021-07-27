package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/just1689/kubernetes-warm-images/model"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var emptyArr = []string{}

func ConnectToKubernetesAPI() *K8sClient {
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

func (cl *K8sClient) WatchEachNamespace(namespaces chan string, nsIgnore []string) chan string {
	result := make(chan string, 256)
	util.FuncForEach(namespaces, func(namespace string) {
		logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("New watcher ns: '%s'", namespace)))
		wi := cl.newWatcher(namespace)
		go watchOneNamespace(namespace, wi, result, nsIgnore)
	})
	return result
}

func watchOneNamespace(namespace string, wi watch.Interface, imgChan chan string, nsIgnore []string) {
	logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("watchOneNamespace namespace: '%s'", namespace)))
	for event := range wi.ResultChan() {
		util.StrArrToChan(getImages(event, nsIgnore), imgChan)
	}

	//allPods := eventsToPods(wi.ResultChan())
	//filteredPods := applyControllerFilters(allPods) //TODO: send Controller filters?
	//images := podsToImages(filteredPods)
	//publishImages(images, imgChan)
	//TODO: rethink
	/*
		--> Transform event into strongly typed object
		--> Process through list of filters?
		--> return image
	*/
}

func eventsToPods(in <-chan watch.Event) chan *v1.Pod {
	result := make(chan *v1.Pod)
	go func() {
		for event := range in {
			if p, ok := event.Object.(*v1.Pod); !ok {
				//TODO: handle !ok
			} else {
				result <- p
			}
		}
		close(result)
	}()
	return result
}

func applyControllerFilters(pods chan *v1.Pod) chan *v1.Pod {
	result := make(chan *v1.Pod)
	go func() {
		for pod := range pods {
			ignore := false
			//
			// TODO: Apply filters
			//
			if !ignore {
				result <- pod
			}

		}
		close(result)
	}()
	return result
}

func podsToImages(pods chan *v1.Pod) chan model.Image {
	result := make(chan model.Image)
	go func() {
		for pod := range pods {
			for _, container := range pod.Spec.Containers {
				result <- model.Image{
					Namespace: pod.ObjectMeta.Namespace,
					PodName:   pod.ObjectMeta.Name,
					Labels:    pod.ObjectMeta.Labels,
					Image:     container.Image,
				}
			}
		}
		close(result)
	}()
	return result
}

func publishImages(images chan model.Image, out chan string) {
	go func() {
		for image := range images {
			b, err := json.Marshal(image)
			if err != nil {
				logrus.Errorln("could not convert Image to json", image, err)
				continue
			}
			out <- string(b)
		}
	}()
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
