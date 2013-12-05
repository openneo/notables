package main

import (
	"fmt"
	rt "github.com/dancannon/gorethink"
	. "github.com/openneo/neopets-notables-go/notables"
	"net/http"
	"net/url"
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
	parsedYear, _ := strconv.ParseInt(m[1], 10, 0)
	parsedMonth, _ := strconv.ParseInt(m[2], 10, 0)
	parsedDay, _ := strconv.ParseInt(m[3], 10, 0)
	year, month, day := int(parsedYear), time.Month(parsedMonth), int(parsedDay)

	now := time.Now().In(timeLocation)
	nowYear, nowMonth, nowDay := now.Year(), now.Month(), now.Day()
	if year == nowYear && month == nowMonth && day == nowDay {
		// It's today! How exciting! We'll have a new notable in approximately
		// five minutes, so cache until then.
		writeExpiresIn(w, time.Duration(5)*time.Minute, now)
	} else if year < nowYear || (year == nowYear && month < nowMonth) ||
		(year == nowYear && month == nowMonth && day < nowDay) {
		// This day has already passed. It's a permanent resource, but let's
		// not get crazy. Cache it for 24 hours.
		writeExpiresIn(w, time.Duration(24)*time.Hour, now)
	} else {
		// It's a future date. We're not going to have anything to say about it
		// until that day has come.
		expiry := time.Date(year, month, day, 0, 0, 0, 0, timeLocation)
		writeExpiresAt(w, expiry, now)
	}

	serveNotablesJSONFromDate(w, r, s, int(year), time.Month(month), int(day))
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

	expiry := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0,
		timeLocation)
	writeExpiresAt(w, expiry, now)

	day := now.Add(-time.Duration(dayAgoCount) * time.Hour * 24)
	newPath := fmt.Sprintf("/api/1/days/%d/%d/%d", day.Year(), day.Month(),
		day.Day())
	newURL := url.URL{r.URL.Scheme, r.URL.Opaque, r.URL.User, r.URL.Host,
		newPath, r.URL.RawQuery, r.URL.Fragment}
	http.Redirect(w, r, newURL.String(), 307)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *rt.Session),
	session *rt.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, session)
	}
}
