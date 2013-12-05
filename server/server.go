package main

import (
	"flag"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/openneo/neopets-notables-go/db"
	"log"
	"net/http"
)

func main() {
	db.SetupFlag()
	port := flag.Int("port", 8888, "port on which to run web server")
	memcacheAddr := flag.String("memcache-address", "localhost:11211",
		"address of the memcache server")
	flag.Parse()

	session, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	mc := memcache.New(*memcacheAddr)

	http.HandleFunc("/api/1/days/ago/", handleDayAgo)
	http.HandleFunc("/api/1/days/", makeHandler(handleExactDay, session, mc))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
