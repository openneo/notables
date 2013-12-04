package main

import (
	"flag"
	"fmt"
	"github.com/openneo/neopets-notables-go/source"
	"log"
	"strconv"
)

func main() {
	flag.Parse()
	maxTries, err := strconv.ParseUint(flag.Arg(0), 10, 64)
	if err != nil {
		maxTries = 1
	}
	notable, ok := source.GetNotable(maxTries)
	if !ok {
		log.Fatalf("failed after %d tries\n", maxTries)
	}
	fmt.Println("success! ", notable)
}
