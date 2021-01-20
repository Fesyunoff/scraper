//+build !swipe

// Code generated by Swipe v2.0.0-beta6. DO NOT EDIT.

package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
)

type ClientOption func(*clientOpts)
type clientOpts struct {
	scraperRequestGetAvailabilityClientOption       []http.ClientOption
	scraperRequestGetAvailabilityEndpointMiddleware []endpoint.Middleware
	scraperRequestGetResponceTimeClientOption       []http.ClientOption
	scraperRequestGetResponceTimeEndpointMiddleware []endpoint.Middleware
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

type clientScraperRequest struct {
	getAvailabilityEndpoint endpoint.Endpoint
	getResponceTimeEndpoint endpoint.Endpoint
}

func (c *clientScraperRequest) GetAvailability(ctx context.Context, site string) (string, error) {
	resp, err := c.getAvailabilityEndpoint(ctx, GetAvailabilityRequest{Site: site})
	if err != nil {
		return "", err
	}
	response := resp.(string)
	return response, nil
}

func (c *clientScraperRequest) GetResponceTime(ctx context.Context, limit string) (string, error) {
	resp, err := c.getResponceTimeEndpoint(ctx, GetResponceTimeRequest{Limit: limit})
	if err != nil {
		return "", err
	}
	response := resp.(string)
	return response, nil
}
