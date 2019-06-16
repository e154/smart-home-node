package main

import (
	"os"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/e154/smart-home-node/system/graceful_service"
	"github.com/op/go-logging"
	"github.com/e154/smart-home-node/system/client"
	l "github.com/e154/smart-home-node/system/logging"
)

var (
	log = logging.MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			fmt.Printf(shortVersionBanner, GetHumanVersion())
			return
		default:
			fmt.Printf(verboseVersionBanner, "v2", os.Args[0])
			return
		}
	}

	start()
}

func start() {

	fmt.Printf(shortVersionBanner, "")

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
