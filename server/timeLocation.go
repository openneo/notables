package main

import (
	"log"
	"time"
)

var timeLocation, _ = time.LoadLocation("America/Los_Angeles")

func init() {
	if timeLocation == nil {
		log.Fatalln("can't load time location")
	}
}
