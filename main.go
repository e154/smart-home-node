package main

import (
	"github.com/e154/smart-home-node/settings"
	"github.com/e154/smart-home-node/server"
	"github.com/e154/smart-home-node/cache"
	"log"
	"time"
)

func main() {
	// settings
	st := settings.SettingsPtr()
	st.Init()

	// cache
	cache.Init(int64(st.Cachetime))

	log.Printf("Start node v%s\n", st.AppVresion())

	// rpc server
	sr := server.ServerPtr()
	if err := sr.Start(st.IP, st.Port); err != nil {
		log.Fatal(err.Error())
	}

	for ;; {
		time.Sleep(time.Second)
	}
}