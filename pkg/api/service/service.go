package service

import (
	"context"

	"github.com/fesyunoff/availability/pkg/types"
)

type ScraperRequest interface {
	GetAvailability(ctx context.Context, site string, id string) (msg string, err error)
	GetResponceTime(ctx context.Context, limit string, id string) (responce string, err error)
	GetStatistics(ctx context.Context, hours string, limit string, id string) (responce []types.Stat, err error)
}
