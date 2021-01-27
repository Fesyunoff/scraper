package scraper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fesyunoff/availability/pkg/api/service"
	"github.com/fesyunoff/availability/pkg/storage"
	"github.com/fesyunoff/availability/pkg/types"
)

type Scraper struct {
	db *storage.PostgreScraperStorage
}

var _ service.ScraperRequest = (*Scraper)(nil)

func (s *Scraper) GetAvailability(ctx context.Context, site string, id string) (responce string, err error) {
	userId, _ := strconv.Atoi(id)
	_, err = s.db.ReturnUsersRole(s.db.Conn, userId)
	if err != nil {
		err = errors.New("ERROR: permission denied")
		log.Println(err)
		return
	}

	row, err := s.db.DisplayServiceAvailability(s.db.Conn, site)
	t := time.Unix(row.Date, 0)
	// t_str := t.
	if !row.Responce {
		responce = fmt.Sprintf("%+v: Service %s is not available (ERROR: timeout)", t, site)
	} else {
		responce = fmt.Sprintf("%+v: Service %s return status code: %d at %d ms", t, site, row.StatusCode, row.Duration)
	}
	return
}

func (s *Scraper) GetResponceTime(ctx context.Context, limit string, id string) (responce string, err error) {
	userId, _ := strconv.Atoi(id)
	_, err = s.db.ReturnUsersRole(s.db.Conn, userId)
	if err != nil {
		err = errors.New("ERROR: permission denied")
		log.Println(err)
		return
	}

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

func (s *Scraper) GetStatistics(ctx context.Context, h string, lim string, id string) (responce []types.Stat, err error) {
	userId, _ := strconv.Atoi(id)
	admin, err := s.db.ReturnUsersRole(s.db.Conn, userId)
	if err != nil || !admin {
		err = errors.New("ERROR: permission denied")
		log.Println(err)
		return
	}

	hours, _ := strconv.ParseInt(h, 10, 64)
	limit, _ := strconv.Atoi(lim)
	return s.db.DisplayStatistics(s.db.Conn, s.db.Conf, hours, limit)
}

func NewScraper(db *storage.PostgreScraperStorage) *Scraper {
	return &Scraper{
		db: db,
	}
}
