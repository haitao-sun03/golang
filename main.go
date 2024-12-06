package main

import (
	"github.com/haitao-sun03/go/controllers"

	"github.com/haitao-sun03/go/config"
	"github.com/haitao-sun03/go/routers"
)

func main() {
	config.Init()
	controllers.InitAccountPathInController()
	// go event.ListenEvents()
	r := routers.Router()
	r.Run(":9999")
}
