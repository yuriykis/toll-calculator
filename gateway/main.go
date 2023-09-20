package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "server listen address")
	flag.Parse()
	http.HandleFunc("/invoice", handleGetInvoice)
	logrus.Infof("gateway HTTP server started: %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handleGetInvoice(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
