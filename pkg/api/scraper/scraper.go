package scraper

import (
	"context"
	"fmt"
	"time"

	"github.com/fesyunoff/availability/pkg/api/service"
	"github.com/fesyunoff/availability/pkg/storage"
)

type Scraper struct {
	db *storage.SQLiteScraperStorage
}

var _ service.ScraperRequest = (*Scraper)(nil)

func (s *Scraper) GetAvailability(ctx context.Context, site string) (responce string, err error) {
	row, err := s.db.DisplayServiceAvailability(s.db.Conn, site)
	t := time.Unix(row.Date, 0)
	// t_str := t.
	if row.Responce == 0 {
		responce = fmt.Sprintf("%+v: Service %s is not available (ERROR: timeout)", t, site)
	} else {
		responce = fmt.Sprintf("%+v: Service %s return status code: %d at %d ms", t, site, row.StatusCode, row.Duration)
	}
	return
}

func (s *Scraper) GetResponceTime(ctx context.Context, limit string) (responce string, err error) {
	switch limit {
	case "min":
		min := true

		row, err := s.db.DisplayServiceResponceTime(s.db.Conn, min)
		t := time.Unix(row.Date, 0)
		responce = fmt.Sprintf("%+v: Service %s with status code: %d has MIN responce time %d ms", t, row.Service, row.StatusCode, row.Duration)
		return responce, err

	case "max":
		min := false

		row, err := s.db.DisplayServiceResponceTime(s.db.Conn, min)
		t := time.Unix(row.Date, 0)
		responce = fmt.Sprintf("%+v: Service %s with status code: %d has MAX responce time %d ms", t, row.Service, row.StatusCode, row.Duration)
		return responce, err

	default:

		responce = "ERROR: Uncorrect request"
		return responce, nil
	}
}

func NewScraper(db *storage.SQLiteScraperStorage) *Scraper {
	return &Scraper{
		db: db,
	}
}
