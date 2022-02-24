package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func main() {
	routerConfiguration := readRouterConfiguration()
	var r *gin.Engine
	if routerConfiguration.Setting.AccessLog {
		r = gin.Default()
	} else {
		r = gin.Default()
	}
	for _, path := range routerConfiguration.paths {
		if strings.HasPrefix(path, "GET") {
			r.GET(routerConfiguration.Configuration[path][0].Request.Url, routerConfiguration.generateGetHandler(path))
		}
	}
	r.NoRoute(func(c *gin.Context) {
		msg := fmt.Sprintf("%s %s page not found", c.Request.Method, c.Request.RequestURI)
		log.Println(msg)
		c.String(404, msg)
	})
	r.Run(":8080")
}
