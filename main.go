package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	r := gin.New()
	routerConfiguration := readRouterConfiguration()

	for _, path := range routerConfiguration.paths {
		if strings.HasPrefix(path, "GET") {
			r.GET(routerConfiguration.Configuration[path][0].Request.Url, routerConfiguration.generateHandler(path))
		}
	}
	r.NoRoute(func(c *gin.Context) {
		msg := fmt.Sprintf("%s %s page not found", c.Request.Method, c.Request.RequestURI)
		log.Println(msg)
		c.String(404, msg)
	})
	r.Run()
}

func (routerConfig RouterConfiguration) generateHandler(path string) func(context *gin.Context) {
	var bodyFile string
	var status int
	var headers map[string]string
	delay := -1
	var requestParaMapping RequestParamMapping
	configurations := routerConfig.Configuration[path]
	requestParaMapping = routerConfig.generateRequestParamMapping(path)

	return func(context *gin.Context) {
		if len(configurations) >= 1 && context.Request.Method == "GET" {
			//can't decide util request comes
			matchedMapping, error := requestParaMapping.getMappingWithRequestQuery(context.Request.URL.Query())
			if error != nil {
				msg := context.Request.RequestURI + "  not found"
				context.String(404, msg)
				log.Println(msg)
				return
			} else {
				bodyFile = matchedMapping.Response.BodyFileName
				status = matchedMapping.Response.Status
				headers = matchedMapping.Response.Headers
				delay = matchedMapping.Response.FixDelay
			}
		}
		delay = routerConfig.Setting.delayDistribution.getDelay(delay)
		content, err := ioutil.ReadFile(filepath.Join("stub/__files", bodyFile))
		if err != nil {
			msg := context.Request.RequestURI + " file not found: " + bodyFile
			log.Println(msg)
			sendResponse(delay, context.String, 500, msg)
			return
		} else {
			var data map[string]interface{}
			json.Unmarshal(content, &data)
			for header, value := range headers {
				context.Header(header, value)
			}
			if data == nil {
				sendResponse(delay, context.String, status, string(content))
			} else {
				sendResponse(delay, context.String, status, string(content))
			}
		}
	}
}

func sendResponse(delay int, fn func(code int, format string, values ...interface{}), status int, msg string) {
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	fn(status, msg)
}
