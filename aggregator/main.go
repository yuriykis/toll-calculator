package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yuriykis/tolling/types"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		store          = makeStore()
		svc            = NewInvoiceAggregator(store)
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, svc))
	}()

	time.Sleep(time.Second)

	// c, err := client.NewGRPCClient(*grpcListenAddr)
	// if err != nil {
	// 	log.Fatal(err)

	// }
	// if err := c.Aggregate(context.Background(), &types.AggregateRequest{
	// 	ObuId: 1,
	// 	Value: 58.44,
	// 	Unix:  time.Now().Unix(),
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	log.Fatal(makeHTTPTransport(httpListenAddr, svc))
}

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatal("Unknown store type")
		return nil
	}
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("Starting gRPC transport")
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)

	// register our gRPC server implementation with the gRPC package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))

	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	var (
		aggMetricsHandler = newHTTPMetricHandler("aggregate")
		invMetricsHandler = newHTTPMetricHandler("invoice")
	)

	http.HandleFunc("/aggregate", aggMetricsHandler.instrument(handleAggregate(svc)))
	http.HandleFunc("/invoice", invMetricsHandler.instrument(handleGetInvoice(svc)))
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	fmt.Println("Starting HTTP transport")
	return http.ListenAndServe(listenAddr, nil)
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
