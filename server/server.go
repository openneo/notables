package main

import (
	"flag"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/openneo/neopets-notables-go/db"
	. "github.com/openneo/neopets-notables-go/notables"
	"log"
	"time"
)

var timeLocation, _ = time.LoadLocation("America/Los_Angeles")

func getNotablesFromDate(session *r.Session, year int, month time.Month, day int) (
	*r.ResultRows, error) {
	// "The month, day, hour, min, sec, and nsec values may be outside their
	//  usual ranges and will be normalized during the conversion. For example,
	//  October 32 converts to November 1." (http://golang.org/pkg/time/#Date)
	// We use Go's time package, rather than RethinkDB's, because Go will
	// adjust by DST by location, whereas RethinkDB requires an hour offset.
	startIncl := time.Date(year, month, day, 0, 0, 0, 0, timeLocation)
	endExcl := time.Date(year, month, day+1, 0, 0, 0, 0, timeLocation)

	return r.Table("notables").Filter(
		r.Row.Field("observed").During(startIncl, endExcl)).Run(session)
}

func printNotables(rows *r.ResultRows) error {
	for rows.Next() {
		var notable Notable
		err := rows.Scan(&notable)
		if err != nil {
			log.Println("can't read row: ", err)
		}
		fmt.Println(notable)
	}

	return nil
}

func init() {
	if timeLocation == nil {
		log.Fatalln("can't load time location")
	}
}

func main() {
	db.SetupFlag()
	flag.Parse()

	session, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	now := time.Now().In(timeLocation)
	notables, err := getNotablesFromDate(session, now.Year(), now.Month(), now.Day())
	if err != nil {
		log.Fatalln(err)
	}

	err = printNotables(notables)
	if err != nil {
		log.Fatalln(err)
	}
}
