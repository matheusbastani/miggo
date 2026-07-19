package settings

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Driver string

const (
	DriverPostgres Driver = "postgres"
)

func newDriver(driver Driver, url string) (*sql.DB, error) {
	switch driver {
	case DriverPostgres:
		return sql.Open("postgres", url)
	case "":
		return nil, fmt.Errorf("driver not specified")
	default:
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}
}
