package main

import (
	"encoding/json"
	"fmt"
	rt "github.com/dancannon/gorethink"
	. "github.com/openneo/neopets-notables-go/notables"
	"net/http"
	"time"
)

func renderJSON(w http.ResponseWriter, r *http.Request, b []byte) {
	callback := r.FormValue("callback")
	if callback == "" {
		fmt.Fprintf(w, "%s", b)
	} else {
		fmt.Fprintf(w, "%s(%s);", callback, b)
	}
}

func renderNotablesJSON(rows *rt.ResultRows) ([]byte, error) {
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

	response := notablesResponse{notables}
	b, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func serveNotablesJSONFromDate(w http.ResponseWriter, r *http.Request, s *rt.Session, year int, month time.Month, day int) {
	notables, err := getNotablesFromDate(s, year, month, day)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := renderNotablesJSON(notables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, r, b)
}

func writeExpiryHeaders(w http.ResponseWriter, now time.Time, expiry time.Time) {
	secondsUntilExpiry := int(expiry.Sub(now).Seconds())
	if secondsUntilExpiry < 0 {
		secondsUntilExpiry = 0
	}

	w.Header().Add("cache-control",
		fmt.Sprintf("public, max-age=%d", secondsUntilExpiry))
	w.Header().Add("expires", expiry.Format(time.RFC1123))
}
