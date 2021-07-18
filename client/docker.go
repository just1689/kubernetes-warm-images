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

func ConnectToDocker() func(next string) {
	cli, err := eclient.NewClient(DockerSocket, "v1.22", nil, nil)
	if err != nil {
		logrus.Panicln(util.LogPrepend(3, fmt.Sprintf("could not connect to docker socket at '%s'. Exiting", DockerSocket)), err)
	}
	return func(image string) {
		image = util.PrependImage(image)
		logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("pull image: '%s'", image)))
		if _, err = cli.ImagePull(context.Background(), image, types.ImagePullOptions{}); err != nil {
			logrus.Errorln(util.LogPrepend(3, fmt.Sprintf("pull image: '%s' ~ FAIL", image)), err)
			return
		}
		logrus.Infoln(util.LogPrepend(3, fmt.Sprintf("pull image: '%s' ~ PASS", image)))
	}
}
