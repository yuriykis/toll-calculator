package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/aggregator/client"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "server listen address")
	aggregatorServiceAddr := flag.String("aggregatorServiceAddr", "http://localhost:3000", "aggregator service address")
	flag.Parse()
	client := client.NewHTTPClient(*aggregatorServiceAddr)
	ih := NewInvoiceHandler(client)
	http.HandleFunc("/invoice", makeAPIFunc(ih.handleGetInvoice))
	logrus.Infof("gateway HTTP server started: %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	obuIDStr := r.URL.Query().Get("obu")
	if obuIDStr == "" {
		return writeJSON(w, http.StatusBadRequest, map[string]string{"error": "obu is required"})
	}
	obuID, err := strconv.Atoi(obuIDStr)
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu"})
	}
	inv, err := h.client.GetInvoice(context.Background(), obuID)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"method": r.Method,
				"uri":    r.RequestURI,
				"took":   time.Since(start),
			}).Info("request processed")
		}(time.Now())

		if err := fn(w, r); err != nil {
			writeJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()},
			)
		}
	}
}
