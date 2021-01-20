//+build !swipe

// Code generated by Swipe v2.0.0-beta6. DO NOT EDIT.

package transport

import (
	http2 "net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
)

type httpError struct {
	code int
}

func (e *httpError) Error() string {
	return http2.StatusText(e.code)
}
func (e *httpError) StatusCode() int {
	return e.code
}
func ErrorDecode(code int) (err error) {
	switch code {
	default:
		err = &httpError{code: code}
	}
	return
}

func middlewareChain(middlewares []endpoint.Middleware) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		if len(middlewares) == 0 {
			return next
		}
		outer := middlewares[0]
		others := middlewares[1:]
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}
func GenericServerOptions(v ...http.ServerOption) ServerOption {
	return func(o *serverOpts) { o.genericServerOption = v }
}

func GenericServerEndpointMiddlewares(v ...endpoint.Middleware) ServerOption {
	return func(o *serverOpts) { o.genericEndpointMiddleware = v }
}

type ServerOption func(*serverOpts)
type serverOpts struct {
	genericServerOption                             []http.ServerOption
	genericEndpointMiddleware                       []endpoint.Middleware
	scraperRequestGetAvailabilityServerOption       []http.ServerOption
	scraperRequestGetAvailabilityEndpointMiddleware []endpoint.Middleware
	scraperRequestGetResponceTimeServerOption       []http.ServerOption
	scraperRequestGetResponceTimeEndpointMiddleware []endpoint.Middleware
}

func ScraperRequestGetAvailabilityServerOptions(opt ...http.ServerOption) ServerOption {
	return func(c *serverOpts) { c.scraperRequestGetAvailabilityServerOption = opt }
}

func ScraperRequestGetAvailabilityServerEndpointMiddlewares(opt ...endpoint.Middleware) ServerOption {
	return func(c *serverOpts) { c.scraperRequestGetAvailabilityEndpointMiddleware = opt }
}

func ScraperRequestGetResponceTimeServerOptions(opt ...http.ServerOption) ServerOption {
	return func(c *serverOpts) { c.scraperRequestGetResponceTimeServerOption = opt }
}

func ScraperRequestGetResponceTimeServerEndpointMiddlewares(opt ...endpoint.Middleware) ServerOption {
	return func(c *serverOpts) { c.scraperRequestGetResponceTimeEndpointMiddleware = opt }
}
