package main

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	rt "github.com/dancannon/gorethink"
	. "github.com/openneo/neopets-notables-go/notables"
	"log"
	"net/http"
	"time"
)

func serveJSON(w http.ResponseWriter, r *http.Request, b []byte) {
	callback := r.FormValue("callback")
	if callback == "" {
		fmt.Fprintf(w, "%s", b)
	} else {
		fmt.Fprintf(w, "%s(%s);", callback, b)
	}
}

func renderNotablesJSON(notables []Notable) ([]byte, error) {
	response := notablesResponse{notables}
	return json.Marshal(response)
}

func renderNotableRowsJSON(rows *rt.ResultRows) ([]byte, error) {
	notables := []Notable{}

	for rows.Next() {
		var notable Notable
		err := rows.Scan(&notable)
		if err != nil {
			return nil, err
		}
		notable.Observed = notable.Observed.In(timeLocation)
		notables = append(notables, notable)
	}

	return renderNotablesJSON(notables)
}

func getNotablesJSON(s *rt.Session, year int, month time.Month, day int) ([]byte, error) {
	log.Printf("Getting from database: %d/%d/%d", year, month, day)
	notables, err := getNotablesFromDate(s, year, month, day)
	if err != nil {
		return nil, err
	}

	b, err := renderNotableRowsJSON(notables)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func getCachedNotablesJSON(s *rt.Session, mc *memcache.Client, year int,
	month time.Month, day int) ([]byte, error) {
	log.Printf("Getting from cache: %d/%d/%d", year, month, day)
	cacheKey := fmt.Sprintf("openneo_notables:days:%d:%d:%d:json", year, month, day)
	var b []byte
	item, err := mc.Get(cacheKey)
	if err == nil {
		b = item.Value
	} else {
		// Couldn't fetch from cache; build it by hand, instead.
		b, err = getNotablesJSON(s, year, month, day)
		if err != nil {
			return nil, err
		}
		err = mc.Set(&memcache.Item{Key: cacheKey, Value: b})
		if err != nil {
			// We don't need to crash if the cache set fails, but we do want
			// to have it on record.
			log.Println("cache set error: ", err)
		}
	}
	return b, nil
}

func serveNotablesJSONFromDate(w http.ResponseWriter, r *http.Request,
	s *rt.Session, mc *memcache.Client, year int, month time.Month, day int) {
	// Work out the expiry info and write the HTTP headers.
	now := time.Now().In(timeLocation)
	nowYear, nowMonth, nowDay := now.Year(), now.Month(), now.Day()
	var b []byte
	var err error
	if year == nowYear && month == nowMonth && day == nowDay {
		// It's today! How exciting! We'll have a new notable in approximately
		// five minutes, so HTTP cache until then. Don't Memcache at all, since
		// it'll just expire soon anyway. CONSIDER: If we get heavy traffic,
		// maybe it'll be worth caching for those 5 minutes.
		writeExpiresIn(w, time.Duration(5)*time.Minute, now)
		b, err = getNotablesJSON(s, year, month, day)
	} else if year < nowYear || (year == nowYear && month < nowMonth) ||
		(year == nowYear && month == nowMonth && day < nowDay) {
		// This day has already passed. It's a permanent resource, but let's
		// not get crazy. Cache it for 24 hours on the client, but we'll keep
		// it semi-permanently in memcache because we can clear that at our
		// leisure.
		writeExpiresIn(w, time.Duration(24)*time.Hour, now)
		b, err = getCachedNotablesJSON(s, mc, year, month, day)
	} else {
		// It's a future date. We're not going to have anything to say about it
		// until that day has come, and we definitely don't have any data for
		// it yet.
		expiry := time.Date(year, month, day, 0, 0, 0, 0, timeLocation)
		writeExpiresAt(w, expiry, now)

		// Serve an empty set. Don't even bother checking the cache or db.
		b, err = renderNotablesJSON([]Notable{})
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	serveJSON(w, r, b)
}

func writeExpiresAt(w http.ResponseWriter, expiry time.Time, now time.Time) {
	writeExpiryHeaders(w, expiry, expiry.Sub(now))
}

func writeExpiresIn(w http.ResponseWriter, timeUntilExpiry time.Duration,
	now time.Time) {
	writeExpiryHeaders(w, now.Add(timeUntilExpiry), timeUntilExpiry)
}

func writeExpiryHeaders(w http.ResponseWriter, expiry time.Time,
	timeUntilExpiry time.Duration) {
	secondsUntilExpiry := int(timeUntilExpiry.Seconds())
	if secondsUntilExpiry < 0 {
		secondsUntilExpiry = 0
	}

	w.Header().Add("cache-control",
		fmt.Sprintf("public, max-age=%d", secondsUntilExpiry))
	w.Header().Add("expires", expiry.Format(time.RFC1123))
}
