package storage

import (
	"database/sql"

	"github.com/fesyunoff/availability/pkg/types"
)

type ScraperStorage interface {
	DisplayServiceAvailability(db *sql.DB, site string) (row types.Row, err error)
	DisplayServiceResponceTime(db *sql.DB, min bool) (r types.Row, err error)
	DisplayStatistics(db *sql.DB, c *types.Config, hours int64, limit int) (out []types.Stat, err error)
	ReturnUsersRole(db *sql.DB, id int) (admin bool, err error)
}
