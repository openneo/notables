package db

import (
	"flag"
	r "github.com/dancannon/gorethink"
)

var address, database, authkey *string

func SetupFlag() {
	address = flag.String("address", "localhost:28015", "rethinkdb address")
	database = flag.String("database", "test", "database name")
	authkey = flag.String("authkey", "", "database authentication key")
}

func Connect() (*r.Session, error) {
	return r.Connect(map[string]interface{}{
		"address":  *address,
		"database": *database,
		"authkey":  *authkey,
	})
}
