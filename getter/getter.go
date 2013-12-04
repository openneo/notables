package main

import (
	"flag"
	"fmt"
	r "github.com/dancannon/gorethink"
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

func save(notable Notable, address string, database string, authkey string) error {
	session, err := r.Connect(map[string]interface{}{
		"address":  address,
		"database": database,
		"authkey":  authkey,
	})
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
	address := flag.String("address", "localhost:28015", "rethinkdb address")
	database := flag.String("database", "test", "database name")
	authkey := flag.String("authkey", "", "database authentication key")
	flag.Parse()

	maxTries := parseMaxTries(flag.Arg(0))
	notable, ok := source.GetNotable(maxTries)
	if !ok {
		log.Fatalf("failed after %d tries\n", maxTries)
	}
	fmt.Println("got notable ", notable)

	err := save(notable, *address, *database, *authkey)
	if err != nil {
		log.Fatalln("save failed: ", err.Error())
	}
	fmt.Println("saved")
}
