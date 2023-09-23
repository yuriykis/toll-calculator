package aggtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"github.com/yuriykis/tolling/go-kit-example/aggsvc/aggendpoint"
	"github.com/yuriykis/tolling/go-kit-example/aggsvc/aggservice"
	"golang.org/x/time/rate"
)

func NewHTTPClient(instance string, logger log.Logger) (aggservice.Service, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// global client middlewares
	var options []httptransport.ClientOption

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var aggregateEndpoint endpoint.Endpoint
	{
		aggregateEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/aggregate"),
			encodeHTTPGenericRequest,
			decodeHTTPAggregateResponse,
			options...,
		).Endpoint()
		aggregateEndpoint = limiter(aggregateEndpoint)
		aggregateEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{
				Name:    "Aggregate",
				Timeout: 30 * time.Second,
			}),
		)(
			aggregateEndpoint,
		)
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var calculateEndpoint endpoint.Endpoint
	{
		calculateEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/invoice"),
			encodeHTTPGenericRequest,
			decodeHTTPAggregateResponse,
			options...,
		).Endpoint()
		calculateEndpoint = limiter(calculateEndpoint)
		calculateEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{
				Name:    "Calculate",
				Timeout: 10 * time.Second,
			}),
		)(
			calculateEndpoint,
		)
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return aggendpoint.Set{
		AggregateEndpoint: aggregateEndpoint,
		CalculateEndpoint: calculateEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func encodeHTTPGenericRequest(
	_ context.Context,
	r *http.Request,
	request interface{},
) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeHTTPAggregateResponse(
	_ context.Context,
	r *http.Response,
) (interface{}, error) {
	var response aggendpoint.AggregateResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	return response, err
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	fmt.Printf("encodeError %v\n", err)
}

func NewHTTPHandler(endpoints aggendpoint.Set, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	mux := http.NewServeMux()
	mux.Handle("/aggregate", httptransport.NewServer(
		endpoints.AggregateEndpoint,
		decodeHTTPAggregateRequest,
		encodeHTTPAggregateResponse,
		options...,
	))
	mux.Handle("/invoice", httptransport.NewServer(
		endpoints.CalculateEndpoint,
		decodeHTTPCalculateRequest,
		encodeHTTPCalculateResponse,
		options...,
	))
	return mux
}

func decodeHTTPAggregateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request aggendpoint.AggregateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeHTTPAggregateResponse(
	_ context.Context,
	w http.ResponseWriter,
	response interface{},
) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeHTTPCalculateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request aggendpoint.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeHTTPCalculateResponse(
	_ context.Context,
	w http.ResponseWriter,
	response interface{},
) error {
	return json.NewEncoder(w).Encode(response)
}
