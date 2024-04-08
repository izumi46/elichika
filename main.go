package main

import (
	"elichika/config"
	"elichika/router"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	router.Router(r)

	if strings.HasPrefix(config.Conf.Settings.HostName, "unix/") {
		path := strings.TrimPrefix(config.Conf.Settings.HostName, "unix/")
		fmt.Println("Running on unix socket: ", path)
		_ = r.RunUnix(path)
	} else {
		addr := config.Conf.Settings.HostName + ":" + config.Conf.Settings.Port
		fmt.Println("Running on: ", addr)
		r.Run(addr)
	}
}
