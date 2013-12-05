package main

import (
	"encoding/json"
	"fmt"
	rt "github.com/dancannon/gorethink"
	. "github.com/openneo/neopets-notables-go/notables"
	"net/http"
	"time"
)

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

func renderNotablesJSONFromDate(w http.ResponseWriter, s *rt.Session, year int, month time.Month, day int) {
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
	fmt.Fprintf(w, "%s", b)
}
