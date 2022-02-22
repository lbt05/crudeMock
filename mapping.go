package main

import (
	"errors"
	"fmt"
	"net/url"
)

type RequestParamMapping struct {
	StandAlone        *Mapping
	ParamMapping      map[string]Mapping
	ParamNamesInOrder []string
}

type MappingConfiguration struct {
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
	FixDelay     int               `json:"fixedDelayMilliseconds"`
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
