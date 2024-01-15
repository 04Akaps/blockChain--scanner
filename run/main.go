package main

import (
	"flag"
	"scanner/env"
	"scanner/log"
	"scanner/run/app"
)

var envFlag = flag.String("env", "./env.toml", "env not found")

func main() {
	flag.Parse()
	e := env.NewEnv(*envFlag)
	log.SetLog(e.Log.LogName)
	app.NewApp(e)
}
