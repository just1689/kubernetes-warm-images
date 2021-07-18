package client

import (
	"fmt"
	eclient "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
)

// DockerSocket TODO: add to Helm chart
var DockerSocket = util.StrOr(os.Getenv("DOCKER_SOCKET"), "unix:///var/run/docker.sock")

func PullDockerImageGo(image string, logPrepend func(level int, text string) (result string)) {
	image = util.PrependImage(image)
	logrus.Infoln(logPrepend(3, fmt.Sprintf("PullDockerImageGo(%s)", image)))
	cli, err := eclient.NewClient(DockerSocket, "v1.22", nil, nil)
	if err != nil {
		logrus.Errorln(logPrepend(3, fmt.Sprintf("PullDockerImageGo(%s) ~ FAIL", image)))
		logrus.Errorln(logPrepend(3, err.Error()))
		return
	}
	if _, err = cli.ImagePull(context.Background(), image, types.ImagePullOptions{}); err != nil {
		logrus.Errorln(logPrepend(3, fmt.Sprintf("PullDockerImageGo(%s) ~ FAIL", image)))
		logrus.Errorln(logPrepend(3, err.Error()))
		return
	}
	logrus.Infoln(logPrepend(3, fmt.Sprintf("PullDockerImageGo(%s) ~ PASS", image)))
}
