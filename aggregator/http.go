package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/types"
)

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func newHTTPMetricHandler(reqName string) *HTTPMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &HTTPMetricHandler{
		reqCounter: reqCounter,
		reqLatency: reqLatency,
	}
}

func (h *HTTPMetricHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"method":  r.Method,
				"path":    r.URL.Path,
				"latenct": latency,
				"took":    time.Since(start).Seconds(),
			}).Info("request")
			h.reqLatency.Observe(latency)
		}(time.Now())
		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(
				w,
				http.StatusMethodNotAllowed,
				map[string]string{"error": "method not allowed"},
			)
		}
		// value, ok := r.URL.Query()["obu"]
		// if !ok || len(value[0]) < 1 {
		// 	writeJSON(
		// 		w,
		// 		http.StatusBadRequest,
		// 		map[string]string{"error": "missing obu"},
		// 	)
		// 	return
		// }
		params := r.URL.Query()
		obu := params.Get("obu")
		if obu == "" {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "missing obu"},
			)
			return
		}
		obuID, err := strconv.Atoi(obu)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()},
			)
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}
func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(
				w,
				http.StatusMethodNotAllowed,
				map[string]string{"error": "method not allowed"},
			)
		}
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()},
			)
			return
		}
	}
}
