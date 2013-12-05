package main

import (
	"flag"
	"fmt"
	"github.com/openneo/neopets-notables-go/db"
	"log"
	"net/http"
)

func main() {
	db.SetupFlag()
	port := flag.Int("port", 8888, "port on which to run web server")
	flag.Parse()

	session, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/api/1/days/ago/", makeHandler(handleDayAgo, session))
	http.HandleFunc("/api/1/days/", makeHandler(handleExactDay, session))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
