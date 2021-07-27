package model

type Image struct {
	Namespace string            `json:"namespace"`
	PodName   string            `json:"podName"`
	Labels    map[string]string `json:"labels"`
	Image     string            `json:"image"`
}
