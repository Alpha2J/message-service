package main

import (
	"message-service/internal/app"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)

	app.InitScheduler()
	app.InitHttp()
}
