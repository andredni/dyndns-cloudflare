package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Configuration struct {
	Records []string
}

var config Configuration

func loadConfig(extraFile string) Configuration {
	jsonData, err := ioutil.ReadFile(extraFile)

	if err != nil {
		fmt.Println(err)
		return config
	}

	json.Unmarshal(jsonData, &config)
	log.Println("Load extra zone file " + extraFile + ": " + strings.Join(config.Records, ", "))
	return config
}
