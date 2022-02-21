package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	r := gin.New()
	routerConfiguration := readConfigurations()

	for _, path := range routerConfiguration.paths {
		if strings.HasPrefix(path, "GET") {
			r.GET(routerConfiguration.Configuration[path][0].Request.Url, routerConfiguration.generateHandler(path))
		}
	}
	r.NoRoute(func(c *gin.Context) {
		msg := c.Request.RequestURI + " page not found"
		log.Println(msg)
		c.String(404, msg)
	})
	r.Run()
}

func (routerConfig RouterConfiguration) generateHandler(path string) func(context *gin.Context) {
	var bodyFile string
	var status int
	var headers map[string]string
	var requestParaMapping RequestParamMapping
	configurations := routerConfig.Configuration[path]
	if len(configurations) == 1 {
		bodyFile = configurations[0].Response.BodyFileName
		status = configurations[0].Response.Status
		headers = configurations[0].Response.Headers
	} else {
		requestParaMapping = routerConfig.generateRequestParamMapping(path)
	}

	return func(context *gin.Context) {
		if len(configurations) > 1 && context.Request.Method == "GET" {
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
			}
		}

		content, err := ioutil.ReadFile(filepath.Join("stub/__files", bodyFile))
		if err != nil {
			msg := context.Request.RequestURI + " file not found: " + bodyFile
			context.String(500, msg)
			log.Println(msg)
			return
		} else {
			var data map[string]interface{}
			json.Unmarshal(content, &data)
			for header, value := range headers {
				context.Header(header, value)
			}
			if data == nil {
				context.String(status, string(content))
			} else {
				context.JSON(status, data)
			}
		}
	}
}

func (mapping RequestParamMapping) getMappingWithRequestQuery(query url.Values) (*Mapping, error) {
	queryKey := ""
	for _, paramName := range mapping.ParamNamesInOrder {
		queryValue := query.Get(paramName)
		if queryValue != "" {
			queryKey += fmt.Sprintf("%s=%v&", paramName, queryValue)
		}
	}
	if matchedMapping, ok := mapping.ParamMapping[queryKey]; ok {
		return &matchedMapping, nil
	} else if mapping.StandAlone != nil {
		return mapping.StandAlone, nil
	}

	return nil, errors.New("mapping not found")
}

func readConfigurations() RouterConfiguration {
	mappingDir := "stub/mappings"
	files, err := ioutil.ReadDir(mappingDir)
	var mappings []Mapping
	if err != nil {
		log.Panic(err)
	} else {
		for _, file := range files {
			if !file.IsDir() {
				jsonFile, ex := ioutil.ReadFile(filepath.Join(mappingDir, file.Name()))

				if ex != nil {
					log.Panic(err)
				} else {
					conf := Configuration{}
					json.Unmarshal(jsonFile, &conf)
					for _, mapping := range conf.Mapping {
						mappings = append(mappings, mapping)
					}
				}
			}
		}
	}
	confMappings := make(map[string][]Mapping)
	for _, configuration := range mappings {
		confMappings[configuration.Request.Method+configuration.Request.Url] = append(confMappings[configuration.Request.Method+configuration.Request.Url], configuration)
	}
	var routerConfiguration RouterConfiguration
	routerConfiguration.Configuration = confMappings
	routerConfiguration.paths = getMappingKeys(confMappings)
	return routerConfiguration
}

func getMapKeys(myMap map[string]string) []string {
	keys := make([]string, len(myMap))

	i := 0
	for k := range myMap {
		keys[i] = k
		i++
	}
	return keys
}
func getMappingKeys(myMap map[string][]Mapping) []string {
	keys := make([]string, len(myMap))

	i := 0
	for k := range myMap {
		keys[i] = k
		i++
	}
	return keys
}

func (routerConfig RouterConfiguration) generateRequestParamMapping(path string) RequestParamMapping {
	var result RequestParamMapping
	paramMapping := make(map[string]Mapping)
	paramSet := make(map[string]string)
	configurations := routerConfig.Configuration[path]
	for _, conf := range configurations {
		if conf.Request.Params == nil {
			result.StandAlone = &conf
		}
		for param := range conf.Request.Params {
			paramSet[param] = param
		}
	}
	params := getMapKeys(paramSet)
	sort.Strings(params)
	result.ParamNamesInOrder = params
	for _, conf := range configurations {
		paramKey := ""
		for _, param := range params {
			if val, ok := conf.Request.Params[param]; ok {
				paramKey += fmt.Sprintf("%s=%v&", param, val)
			}
		}
		paramMapping[paramKey] = conf
	}
	result.ParamMapping = paramMapping
	return result
}

type RequestParamMapping struct {
	StandAlone        *Mapping
	ParamMapping      map[string]Mapping
	ParamNamesInOrder []string
}

type RouterConfiguration struct {
	Configuration map[string][]Mapping
	paths         []string
}

type Configuration struct {
	Mapping []Mapping `json:"mappings"`
}

type Mapping struct {
	Request  RequestConf  `json:"request"`
	Response ResponseConf `json:"response"`
}

type RequestConf struct {
	Method string                 `json:"method"`
	Url    string                 `json:"url"`
	Params map[string]interface{} `json:"params"`
}

type ResponseConf struct {
	Status       int               `json:"status"`
	BodyFileName string            `json:"bodyFileName"`
	Headers      map[string]string `json:"headers"`
}
