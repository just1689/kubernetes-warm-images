package main

import (
	"github.com/just1689/kubernetes-warm-images/apps"
	"github.com/sirupsen/logrus"
	"os"
)

type AppType string

var (
	AgentType      AppType = "AGENT"
	ControllerType AppType = "CONTROLLER"
	AppMap                 = map[AppType]func(){
		AgentType:      apps.RunAgent,
		ControllerType: apps.RunController,
	}
)

func main() {
	logrus.Infoln("Starting up...")
	startApp()
	<-make(chan os.Signal, 1)
	logrus.Infoln("Signal received. Shutting down")
	os.Exit(0)
}

func startApp() {
	var appToStart = AppType(os.Getenv("APP"))
	app, found := AppMap[appToStart]
	if !found {
		logrus.Panicln("no APP env variable to run from among ", getKeysAsStr(AppMap))
	}
	go app()
}

func getKeysAsStr(appMap map[AppType]func()) string {
	result := ""
	for name, _ := range appMap {
		if result != "" {
			result += ", "
		}
		result += string(name)
	}
	return result
}
