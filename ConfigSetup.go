package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// refer to README.md for setting up Google and Facebook logins

type Config struct {
	ClientID_FB     string
	ClientSecret_FB string
	ClientID        string
	ClientSecret    string
	ChatPrivateKey  string
	SSLPrivateKey   string
	SSLCert         string
	Protocol        string
	Host            string
	Port            string
	VideoProtocol   string
	VideoHost       string
	VideoPort       string
}

var ConfigManager struct {
	_conf Config
}

func (conf *Config) url() string {
	return conf.Protocol + "://" + conf.Host + ":" + conf.Port
}

func (cfg *Config) load(file string) {
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, cfg)
}
