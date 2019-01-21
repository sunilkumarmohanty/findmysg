package main

import (
	"log"

	"github.com/sunilkumarmohanty/findmysg/crawler"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot initialize logger. Error : %v", err)
	}
}

func main() {
	logger.Info("Starting application")
	crawler.Crawl()
}
