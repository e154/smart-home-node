package main

import (
	"fmt"
	"github.com/e154/smart-home-node/system/client"
	"github.com/e154/smart-home-node/system/graceful_service"
	l "github.com/e154/smart-home-node/system/logging"
	"github.com/e154/smart-home-node/version"
	"github.com/op/go-logging"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	log = logging.MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			fmt.Printf(version.ShortVersionBanner, version.GetHumanVersion())
			return
		default:
			fmt.Printf(version.VerboseVersionBanner, "v2", os.Args[0])
			return
		}
	}

	start()
}

func start() {

	fmt.Printf(version.ShortVersionBanner, "")

	container := BuildContainer()
	container.Invoke(func(
		graceful *graceful_service.GracefulService,
		lx *logrus.Logger,
		client *client.Client) {

		l.Initialize(lx)
		client.Connect()

		graceful.Wait()
	})
}
