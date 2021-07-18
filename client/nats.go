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
	p.nc.Subscribe(subjectNameImages, func(msg *nats.Msg) { handler(string(msg.Data)) })
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

	username, password := fileStrOrPanic(usernamePath), fileStrOrPanic(passwordPath)
	cl.nc, err = nats.Connect(natsAddr, nats.UserInfo(username, password))
	if err != nil {
		logrus.Fatalln(util.LogPrepend(3, fmt.Sprint("could not connect to NATs server with NATS URL: ", natsAddr, err)))
	}
	return &cl

}
