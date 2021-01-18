package main

import (
	"context"
	"fmt"
	"log"
	adapterIn "ms-go/adapter/in"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)

	// create the handlers
	pc := adapterIn.NewProductsHTTP(l)

	// create a new mux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/products/", pc)
	sm.Handle("/products", pc)

	// create a new server and tune some configurations
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go listenAndServe(s, l)

	handleShutdown(s, l)
}

func listenAndServe(s *http.Server, l *log.Logger) {
	fmt.Printf("Server listening @ %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		l.Fatal(err)
	}
}

func handleShutdown(s *http.Server, l *log.Logger) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received signal. Graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
