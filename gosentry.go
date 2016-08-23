package main

import (
	log "github.com/cihub/seelog"
	"github.com/grindlemire/GoSentry/start"
)

func main() {
	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")

	if err != nil {
		log.Warn("Failed to load logger config ", err)
	} else {
		log.ReplaceLogger(logger)
	}

	start.Run()

}
