/*
Package postgresql provides the library to communicate to postgresql
*/
package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

// Connect initializes postgresql DB
func Connect(datasource string, maxactive, maxidle int) (*sqlx.DB, error) {
	db := sqlx.MustOpen("postgres", datasource)
	db.SetMaxOpenConns(maxactive)
	db.SetMaxIdleConns(maxidle)
	err := db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %s err: %s", datasource, err)
	}
	return db, nil
}
