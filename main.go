package main

import (
	"net/http"
	"time"

	"github.com/haitao-sun03/go/controllers"
	"github.com/haitao-sun03/go/event"

	"github.com/haitao-sun03/go/config"
	"github.com/haitao-sun03/go/routers"
)

func main() {
	defer config.RoutinePool.ReleaseTimeout(5 * time.Second)
	defer config.RedisClient.Close()

	config.Init()
	controllers.InitAccountPathInController()
	// go event.ListenEvents()
	config.RoutinePool.Submit(func() {
		event.ListenEvents()
	})

	config.RoutinePool.Submit(func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			panic(err)
		}
	})
	r := routers.Router()
	r.Run(":9999")
}
