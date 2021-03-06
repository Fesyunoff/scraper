//+build !swipe

// Code generated by Swipe v2.0.0-beta6. DO NOT EDIT.

package transport

import (
	"context"

	"github.com/fesyunoff/availability/pkg/types"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
)

type ClientOption func(*clientOpts)
type clientOpts struct {
	scraperRequestGetAvailabilityClientOption       []http.ClientOption
	scraperRequestGetAvailabilityEndpointMiddleware []endpoint.Middleware
	scraperRequestGetResponceTimeClientOption       []http.ClientOption
	scraperRequestGetResponceTimeEndpointMiddleware []endpoint.Middleware
	scraperRequestGetStatisticsClientOption         []http.ClientOption
	scraperRequestGetStatisticsEndpointMiddleware   []endpoint.Middleware
	genericClientOption                             []http.ClientOption
	genericEndpointMiddleware                       []endpoint.Middleware
}

func GenericClientOptions(opt ...http.ClientOption) ClientOption {
	return func(c *clientOpts) { c.genericClientOption = opt }
}

func GenericClientEndpointMiddlewares(opt ...endpoint.Middleware) ClientOption {
	return func(c *clientOpts) { c.genericEndpointMiddleware = opt }
}

func ScraperRequestGetAvailabilityClientOptions(opt ...http.ClientOption) ClientOption {
	return func(c *clientOpts) { c.scraperRequestGetAvailabilityClientOption = opt }
}

func ScraperRequestGetAvailabilityClientEndpointMiddlewares(opt ...endpoint.Middleware) ClientOption {
	return func(c *clientOpts) { c.scraperRequestGetAvailabilityEndpointMiddleware = opt }
}

func ScraperRequestGetResponceTimeClientOptions(opt ...http.ClientOption) ClientOption {
	return func(c *clientOpts) { c.scraperRequestGetResponceTimeClientOption = opt }
}

func ScraperRequestGetResponceTimeClientEndpointMiddlewares(opt ...endpoint.Middleware) ClientOption {
	return func(c *clientOpts) { c.scraperRequestGetResponceTimeEndpointMiddleware = opt }
}

func ScraperRequestGetStatisticsClientOptions(opt ...http.ClientOption) ClientOption {
	return func(c *clientOpts) { c.scraperRequestGetStatisticsClientOption = opt }
}

func ScraperRequestGetStatisticsClientEndpointMiddlewares(opt ...endpoint.Middleware) ClientOption {
	return func(c *clientOpts) { c.scraperRequestGetStatisticsEndpointMiddleware = opt }
}

type clientScraperRequest struct {
	getAvailabilityEndpoint endpoint.Endpoint
	getResponceTimeEndpoint endpoint.Endpoint
	getStatisticsEndpoint   endpoint.Endpoint
}

func (c *clientScraperRequest) GetAvailability(ctx context.Context, site string, id string) (string, error) {
	resp, err := c.getAvailabilityEndpoint(ctx, GetAvailabilityRequest{Site: site, Id: id})
	if err != nil {
		return "", err
	}
	response := resp.(string)
	return response, nil
}

func (c *clientScraperRequest) GetResponceTime(ctx context.Context, limit string, id string) (string, error) {
	resp, err := c.getResponceTimeEndpoint(ctx, GetResponceTimeRequest{Limit: limit, Id: id})
	if err != nil {
		return "", err
	}
	response := resp.(string)
	return response, nil
}

func (c *clientScraperRequest) GetStatistics(ctx context.Context, hours string, limit string, id string) ([]types.Stat, error) {
	resp, err := c.getStatisticsEndpoint(ctx, GetStatisticsRequest{Hours: hours, Limit: limit, Id: id})
	if err != nil {
		return nil, err
	}
	response := resp.([]types.Stat)
	return response, nil
}
