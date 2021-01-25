package service

import (
	"context"

	"github.com/fesyunoff/availability/pkg/types"
)

type ScraperRequest interface {
	GetAvailability(ctx context.Context, site string) (msg string, err error)
	GetResponceTime(ctx context.Context, limit string) (responce string, err error)
	GetStatistics(ctx context.Context, hours string, limit string) (responce []types.Stat, err error)
}
