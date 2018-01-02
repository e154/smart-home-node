package main

import (
	"log"
	"os"
)

var (
	stdlog, errlog *log.Logger
)

func main() {

	args := os.Args
	switch len(args) {
	case 1:
		stdlog.Printf(shortVersionBanner, "")
		ServiceInitialize()

	case 2:
		switch args[1] {
		case "install", "remove", "start", "stop", "status":
			ServiceInitialize()
		default:
			stdlog.Printf(verboseVersionBanner, "", args[0])
		}
	default:
		stdlog.Printf(verboseVersionBanner, "", args[0])
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}
