package client

import (
	"fmt"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	subjectNameImages = util.StrOr(os.Getenv("IMAGES"), "images")
	usernamePath      = "/config/nats.username"
	passwordPath      = "/config/nats.password"
)

type PubSub struct {
	nc *nats.Conn
}

func (p *PubSub) Publish(image string) {
	p.nc.Publish(subjectNameImages, []byte(image))
}
func (p *PubSub) Subscribe(handler func(image string)) {
	//TODO: disconnect?
	p.nc.Subscribe(subjectNameImages, func(msg *nats.Msg) {
		handler(string(msg.Data))
	})
}
func (p *PubSub) Close() {
	p.nc.Close()
}

func NewPubSubClient() *PubSub {
	natsAddr := os.Getenv("NATS_ADDR")
	if natsAddr == "" {
		logrus.Panicln("could not find env var NATS_ADDR. Exiting")
	}
	cl := PubSub{}
	var err error

	username, password := getNATsCreds()
	cl.nc, err = nats.Connect(natsAddr, nats.UserInfo(username, password))
	if err != nil {
		logrus.Errorln(util.LogPrepend(3, fmt.Sprint("could not connect to NATs server with NATS URL: ", natsAddr)))
		logrus.Errorln(util.LogPrepend(3, err.Error()))
	}
	return &cl

}

func getNATsCreds() (username, password string) {
	var err error
	if username, err = readFileToString(usernamePath); err != nil || username == "" {
		logrus.Panicln(fmt.Sprintf("no NATs username (%s) found at %s", username, usernamePath))
	}
	if password, err = readFileToString(passwordPath); err != nil || password == "" {
		logrus.Panicln(fmt.Sprintf("no NATs password (%s) found at %s", password, passwordPath))
	}
	return
}
