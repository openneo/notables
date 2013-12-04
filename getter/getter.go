package main

import (
	"flag"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/openneo/neopets-notables-go/db"
	. "github.com/openneo/neopets-notables-go/notables"
	"github.com/openneo/neopets-notables-go/source"
	"log"
	"strconv"
)

func parseMaxTries(maxTriesString string) uint64 {
	maxTries, err := strconv.ParseUint(maxTriesString, 10, 64)
	if err != nil {
		maxTries = 1
	}
	return maxTries
}

func save(notable Notable) error {
	session, err := db.Connect()
	if err != nil {
		return err
	}

	_, err = r.Table("notables").Insert(notable).Run(session)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db.SetupFlag()
	flag.Parse()

	maxTries := parseMaxTries(flag.Arg(0))
	notable, ok := source.GetNotable(maxTries)
	if !ok {
		log.Fatalf("failed after %d tries\n", maxTries)
	}
	fmt.Println("got notable ", notable)

	err := save(notable)
	if err != nil {
		log.Fatalln("save failed: ", err.Error())
	}
	fmt.Println("saved")
}
