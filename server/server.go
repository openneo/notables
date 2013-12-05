package main

import (
	"flag"
	"github.com/openneo/neopets-notables-go/db"
	"log"
	"net/http"
)

func main() {
	db.SetupFlag()
	flag.Parse()

	session, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/api/1/days/ago/", makeHandler(handleDayAgo, session))
	http.HandleFunc("/api/1/days/", makeHandler(handleExactDay, session))
	http.ListenAndServe(":8888", nil)
}
