package service

import (
	"context"
)

type ScraperRequest interface {
	GetAvailability(ctx context.Context, site string) (msg string, err error)
	GetResponceTime(ctx context.Context, limit string) (responce string, err error)
}
