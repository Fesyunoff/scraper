//+build !swipe

// Code generated by Swipe v2.0.0-beta6. DO NOT EDIT.

package transport

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	http2 "net/http"
	"net/url"
	"strings"

	"github.com/fesyunoff/availability/pkg/api/service"
	"github.com/go-kit/kit/transport/http"
	"github.com/pquerna/ffjson/ffjson"
)

func NewClientREST(tgt string, options ...ClientOption) (service.ScraperRequest, error) {
	opts := &clientOpts{}
	c := &clientScraperRequest{}
	for _, o := range options {
		o(opts)
	}
	if strings.HasPrefix(tgt, "[") {
		host, port, err := net.SplitHostPort(tgt)
		if err != nil {
			return nil, err
		}
		tgt = host + ":" + port
	}
	u, err := url.Parse(tgt)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	c.getAvailabilityEndpoint = http.NewClient(
		http2.MethodGet,
		u,
		func(_ context.Context, r *http2.Request, request interface{}) error {
			req, ok := request.(GetAvailabilityRequest)
			if !ok {
				return fmt.Errorf("couldn't assert request as GetAvailabilityRequest, got %T", request)
			}
			r.Method = http2.MethodGet
			r.URL.Path += "/getAvailability"
			q := r.URL.Query()
			siteStr := req.Site
			q.Add("site", siteStr)
			r.URL.RawQuery = q.Encode()
			return nil
		},
		func(_ context.Context, r *http2.Response) (interface{}, error) {
			if statusCode := r.StatusCode; statusCode != http2.StatusOK {
				return nil, ErrorDecode(statusCode)
			}
			var resp string
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			err = ffjson.Unmarshal(b, &resp)
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("couldn't unmarshal body to GetAvailabilityRequest: %s", err)
			}
			return resp, nil
		},
		append(opts.genericClientOption, opts.scraperRequestGetAvailabilityClientOption...)...,
	).Endpoint()
	c.getAvailabilityEndpoint = middlewareChain(append(opts.genericEndpointMiddleware, opts.scraperRequestGetAvailabilityEndpointMiddleware...))(c.getAvailabilityEndpoint)
	c.getResponceTimeEndpoint = http.NewClient(
		http2.MethodGet,
		u,
		func(_ context.Context, r *http2.Request, request interface{}) error {
			req, ok := request.(GetResponceTimeRequest)
			if !ok {
				return fmt.Errorf("couldn't assert request as GetResponceTimeRequest, got %T", request)
			}
			r.Method = http2.MethodGet
			r.URL.Path += "/getResponceTime"
			q := r.URL.Query()
			limitStr := req.Limit
			q.Add("limit", limitStr)
			r.URL.RawQuery = q.Encode()
			return nil
		},
		func(_ context.Context, r *http2.Response) (interface{}, error) {
			if statusCode := r.StatusCode; statusCode != http2.StatusOK {
				return nil, ErrorDecode(statusCode)
			}
			var resp string
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			err = ffjson.Unmarshal(b, &resp)
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("couldn't unmarshal body to GetResponceTimeRequest: %s", err)
			}
			return resp, nil
		},
		append(opts.genericClientOption, opts.scraperRequestGetResponceTimeClientOption...)...,
	).Endpoint()
	c.getResponceTimeEndpoint = middlewareChain(append(opts.genericEndpointMiddleware, opts.scraperRequestGetResponceTimeEndpointMiddleware...))(c.getResponceTimeEndpoint)
	return c, nil
}