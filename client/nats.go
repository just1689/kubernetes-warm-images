package client

import (
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"os"
)

var SubjectNameImages = util.StrOr(os.Getenv("IMAGES"), "images")

type PubSub struct {
	nc *nats.Conn
}

func (p *PubSub) Publish(image string) {
	p.nc.Publish(SubjectNameImages, []byte(image))
}
func (p *PubSub) Subscribe(handler func(image string)) {
	//TODO: disconnect?
	p.nc.Subscribe(SubjectNameImages, func(msg *nats.Msg) {
		handler(string(msg.Data))
	})
}
func (p *PubSub) Close() {
	p.nc.Close()
}

func NewPubSubClient() *PubSub {
	natsAddr := os.Getenv("NATS_ADDR")
	if natsAddr == "" {
		logrus.Fatalln("could not find env var NATS_ADDR. Exiting")
	}
	cl := PubSub{}
	var err error
	cl.nc, err = nats.Connect(natsAddr, nats.UserInfo("nats_client", "sY6GYz5c9W"))
	if err != nil {
		logrus.Errorln("could not connect to NATs server with NATS URL: ", natsAddr)
		logrus.Panicln(err)
	}
	return &cl

}
