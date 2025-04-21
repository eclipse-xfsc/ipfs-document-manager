package main

import (
	"log"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/logr"
)

func createLogger(logLevel string, isDev bool) logr.Logger {
	logger, err := logr.New(logLevel, isDev, nil)
	if err != nil {
		log.Fatalf("cannot initialise logger: %t", err)
	}
	return *logger
}
