package main

import (
	"log"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/config"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
}

func getLoadedConfig() Config {
	var conf Config
	err := config.LoadConfig("IPFSMANAGER", &conf, nil)
	if err == nil {
		err = envconfig.Process("IPFSMANAGER", &conf)
	}
	if err != nil {
		log.Fatalf("service was not loaded: %t", err)
	}
	return conf
}
