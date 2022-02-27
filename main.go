package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

func main() {
	routerConfiguration := readRouterConfiguration()
	var r *gin.Engine
	if routerConfiguration.Setting.AccessLog {
		r = gin.Default()
	} else {
		r = gin.New()
	}
	for _, path := range routerConfiguration.paths {
		method, handler := routerConfiguration.generateHandler(path)
		if method == "GET" {
			r.GET(routerConfiguration.Configuration[path][0].Request.Url, handler)
		} else if method == "POST" {
			r.POST(routerConfiguration.Configuration[path][0].Request.Url, handler)
		} else if method == "PUT" {

		}
	}
	r.NoRoute(func(c *gin.Context) {
		msg := fmt.Sprintf("%s %s page not found", c.Request.Method, c.Request.RequestURI)
		log.Println(msg)
		var bodyJson map[string]interface{}
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err == nil {
			json.Unmarshal(jsonData, &bodyJson)
		}

		log.Println(bodyJson)
		c.String(404, msg)
	})
	r.Run(":8080")
}
