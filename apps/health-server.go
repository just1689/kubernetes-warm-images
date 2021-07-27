package apps

import (
	"encoding/json"
	"github.com/just1689/kubernetes-warm-images/health"
	"github.com/just1689/kubernetes-warm-images/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func startHealthServer() {
	logrus.Infoln(util.LogPrepend(2, "starting health server"))
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		logrus.Panicln("could not find env var LISTEN_ADDR. Exiting")
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		ok := health.GlobalHealth.IsSystemOK()
		msg, _ := json.Marshal(ok)
		if !ok.OK {
			http.Error(writer, string(msg), http.StatusInternalServerError)
			return
		}
		writer.Write(msg)
	})
	http.ListenAndServe(listenAddr, nil)
}
