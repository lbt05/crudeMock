package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
)

type RouterConfiguration struct {
	Configuration map[string][]Mapping
	paths         []string
	Setting       Setting
}

func readRouterConfiguration() RouterConfiguration {
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
					conf := MappingConfiguration{}
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

	adminConfiguration := "stub/__admin/settings.json"
	content, ex := ioutil.ReadFile(adminConfiguration)
	if ex == nil {
		conf := Setting{}
		json.Unmarshal(content, &conf)
		routerConfiguration.Setting = conf
	}
	return routerConfiguration
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
