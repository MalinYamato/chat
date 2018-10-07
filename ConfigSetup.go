package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// refer to README.md for setting up Google and Facebook logins

func LoadConfig(file string) *Config {
	var conf *Config
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &conf)

	return conf
}
