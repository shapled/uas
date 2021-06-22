package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"uas/api"
	"uas/settings"

	"github.com/shapled/pitaya"
)

var yamlConfig = flag.String("config", "", "Specific config yaml file")

func init() {
	flag.StringVar(yamlConfig, "c", "", "Specific config yaml file")
}

func main() {
	flag.Parse()

	if yamlConfig != nil {
		if err := settings.SyncFromConfigFile(*yamlConfig); err != nil {
			logrus.Fatal(err)
		}
	}

	server := pitaya.NewServer()
	server.GET("/app/", api.ListApps, &api.AppListRequest{})
	logrus.Fatal(server.Start(":10086"))
}
