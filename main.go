package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"
)

func main() {
	r := gin.New()
	configurations := readConfigurations()

	for path, mappings := range configurations {
		if strings.HasPrefix(path, "GET") {
			r.GET(mappings[0].Request.Url, generateHandler(mappings))
		}
	}
	r.NoRoute(func(c *gin.Context) {
		msg := c.Request.RequestURI + "  not found"
		log.Println(msg)
		c.String(404, msg)
	})
	r.Run()
}

func generateHandler(configurations []Mapping) func(context *gin.Context) {
	var bodyFile string
	var status int
	var headers map[string]string
	var requestParaMapping RequestParamMapping
	if len(configurations) == 1 {
		bodyFile = configurations[0].Response.BodyFileName
		status = configurations[0].Response.Status
		headers = configurations[0].Response.Headers
	} else {
		requestParaMapping = generateRequestParamMapping(configurations)
	}

	return func(context *gin.Context) {
		if len(configurations) > 1 && context.Request.Method == "GET" {
			//can't decide util request comes
			matchedMapping, error := getMappingWithRequestQuery(context.Request.URL.Query(), requestParaMapping.ParamMapping)
			if error != nil {
				if requestParaMapping.StandAlone != nil {
					bodyFile = requestParaMapping.StandAlone.Response.BodyFileName
					status = requestParaMapping.StandAlone.Response.Status
					headers = requestParaMapping.StandAlone.Response.Headers
				} else {
					msg := context.Request.RequestURI + "  not found"
					context.String(404, msg)
					log.Println(msg)
					return
				}
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

func getMappingWithRequestQuery(query url.Values, paramMapping map[string]map[interface{}]Mapping) (Mapping, error) {
	var matchedMappings []Mapping
	for paramName, value := range paramMapping {
		requestParam := query.Get(paramName)
		if mapping, ok := value[requestParam]; ok {
			matchedMappings = append(matchedMappings, mapping)
		} else {
			break
		}
	}
	if matchedMappings != nil {
		allMatch := true
		for i := 1; i < len(matchedMappings)-1; i++ {
			if matchedMappings[i].Response.BodyFileName != matchedMappings[i+1].Response.BodyFileName {
				allMatch = false
				break
			}
		}
		if allMatch {
			return matchedMappings[0], nil
		}
	}
	return Mapping{}, errors.New("mapping not found")
}

func readConfigurations() map[string][]Mapping {
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
	return confMappings
}

func generateRequestParamMapping(configurations []Mapping) RequestParamMapping {
	var paramMapping map[string]map[interface{}]Mapping
	var result RequestParamMapping
	paramMapping = make(map[string]map[interface{}]Mapping)
	for _, conf := range configurations {
		if conf.Request.Params == nil {
			result.StandAlone = &conf
		}
		for param, value := range conf.Request.Params {
			if paramMapping[param] == nil {
				paramMapping[param] = make(map[interface{}]Mapping)
			}
			paramMapping[param][value] = conf
		}
	}
	result.ParamMapping = paramMapping

	return result
}

type RequestParamMapping struct {
	StandAlone   *Mapping
	ParamMapping map[string]map[interface{}]Mapping
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
