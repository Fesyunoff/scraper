//+build !swipe

// Code generated by Swipe v2.0.0-beta6. DO NOT EDIT.

package transport

import (
	"context"

	"github.com/fesyunoff/availability/pkg/api/service"
	"github.com/go-kit/kit/endpoint"
)

func makeScraperRequestGetAvailabilityEndpoint(s service.ScraperRequest) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAvailabilityRequest)
		result, err := s.GetAvailability(ctx, req.Site, req.Id)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

}

func makeScraperRequestGetResponceTimeEndpoint(s service.ScraperRequest) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetResponceTimeRequest)
		result, err := s.GetResponceTime(ctx, req.Limit, req.Id)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

}

func makeScraperRequestGetStatisticsEndpoint(s service.ScraperRequest) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetStatisticsRequest)
		result, err := s.GetStatistics(ctx, req.Hours, req.Limit, req.Id)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

}

type ScraperRequestEndpointSet struct {
	GetAvailabilityEndpoint endpoint.Endpoint
	GetResponceTimeEndpoint endpoint.Endpoint
	GetStatisticsEndpoint   endpoint.Endpoint
}

func MakeScraperRequestEndpointSet(svc service.ScraperRequest) ScraperRequestEndpointSet {
	return ScraperRequestEndpointSet{
		GetAvailabilityEndpoint: makeScraperRequestGetAvailabilityEndpoint(svc),
		GetResponceTimeEndpoint: makeScraperRequestGetResponceTimeEndpoint(svc),
		GetStatisticsEndpoint:   makeScraperRequestGetStatisticsEndpoint(svc),
	}
}

type GetAvailabilityRequest struct {
	Site string `json:"site"`
	Id   string `json:"id"`
}
type GetResponceTimeRequest struct {
	Limit string `json:"limit"`
	Id    string `json:"id"`
}
type GetStatisticsRequest struct {
	Hours string `json:"hours"`
	Limit string `json:"limit"`
	Id    string `json:"id"`
}
