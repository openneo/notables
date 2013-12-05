package main

import (
	r "github.com/dancannon/gorethink"
	"time"
)

func getNotablesFromDate(session *r.Session, year int, month time.Month, day int) (
	*r.ResultRows, error) {
	// "The month, day, hour, min, sec, and nsec values may be outside their
	//  usual ranges and will be normalized during the conversion. For example,
	//  October 32 converts to November 1." (http://golang.org/pkg/time/#Date)
	// We use Go's time package, rather than RethinkDB's, because Go will
	// adjust by DST by location, whereas RethinkDB requires an hour offset.
	startIncl := time.Date(year, month, day, 0, 0, 0, 0, timeLocation)
	endExcl := time.Date(year, month, day+1, 0, 0, 0, 0, timeLocation)

	// TODO: use index for ordering once the next rethinkdb comes out
	betweenOpts := r.BetweenOpts{Index: "observed"}
	return r.Table("notables").
		Between(startIncl, endExcl, betweenOpts).
		OrderBy(r.Row.Field("observed")).
		Run(session)
}
