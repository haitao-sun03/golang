package main

import (
	"github.com/haitao-sun03/go/config"
	"github.com/haitao-sun03/go/routers"
)

func main() {
	r := routers.Router()
	config.Init()
	r.Run(":9999")
}
