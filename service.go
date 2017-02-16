package main

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/takama/daemon"
	"github.com/e154/smart-home-node/settings"
	"github.com/e154/smart-home-node/server"
	"github.com/e154/smart-home-node/cache"
	"github.com/astaxie/beego"
)

const (
	name        = "smart-home-node"
	description = "Smart Home Node"
)

var dependencies = []string{}

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// settings
	st := settings.SettingsPtr()
	st.Init()

	// cache
	cache.Init(int64(st.Cachetime))

	stdlog.Printf("Start node v%s\n", st.AppVresion())

	// rpc server
	sr := server.ServerPtr()
	if err := sr.Start(st.IP, st.Port); err != nil {
		stdlog.Fatal(err.Error())
	}

	var port int
	port = beego.BConfig.Listen.HTTPPort

	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", port)

			if killSignal == os.Interrupt {
				return "Node was interruped by system signal", nil
			}
			return "Node was killed", nil
		}
	}
}

func ServiceInitialize() {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	stdlog.Println(status)
}
