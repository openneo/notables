package main

import (
	rt "github.com/dancannon/gorethink"
	. "github.com/openneo/neopets-notables-go/notables"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type notablesResponse struct {
	Notables []Notable `json:"notables"`
}

var exactDayPath = regexp.MustCompile("^/api/1/days/([0-9]{4})/(0?[1-9]|1[0-2])/(0?[1-9]|[1-2][0-9]|3[0-1])$")

func handleExactDay(w http.ResponseWriter, r *http.Request, s *rt.Session) {
	m := exactDayPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	// Swallowing errors because the regexp already validated its int-ness and
	// size.
	year, _ := strconv.ParseInt(m[1], 10, 0)
	month, _ := strconv.ParseInt(m[2], 10, 0)
	day, _ := strconv.ParseInt(m[3], 10, 0)

	renderNotablesJSONFromDate(w, s, int(year), time.Month(month), int(day))
}

var dayAgoPath = regexp.MustCompile("^/api/1/days/ago/([0-9]+)$")

func handleDayAgo(w http.ResponseWriter, r *http.Request, s *rt.Session) {
	m := dayAgoPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	dayAgoCount, err := strconv.ParseInt(m[1], 10, 0)
	if err != nil {
		// The regexp already validated its int-ness, so an error represents
		// out of bounds. And we certainly don't have more than int days of
		// history. TODO: yield an empty list instead, like with most far-past
		// dates?
		http.NotFound(w, r)
		return
	}

	now := time.Now().In(timeLocation)
	day := now.Add(-time.Duration(dayAgoCount) * time.Hour * 24)

	renderNotablesJSONFromDate(w, s, day.Year(), day.Month(), day.Day())
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *rt.Session),
	session *rt.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, session)
	}
}
