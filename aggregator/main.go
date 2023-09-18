package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yuriykis/tolling/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "server listen address")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("Starting HTTP transport")
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
