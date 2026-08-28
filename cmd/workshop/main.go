package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"workshopnotice/internal/flow"
	"workshopnotice/internal/httpapi"
	"workshopnotice/internal/store"
)

func main() {
	path := flag.String("db", "workshopnotice.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	if err := Run(*path, *addr); err != nil {
		log.Fatal(err)
	}
}

func Run(path, addr string) error {
	db, err := store.Open(path)
	if err != nil {
		return err
	}
	defer db.Close()
	service := flow.NewService(db, func() string { return "2026-01-01T00:00:00Z" })
	handler := httpapi.New(service)
	server := &http.Server{Addr: addr, Handler: httpapi.RegisterMux(handler)}
	fmt.Printf("workshop safety notice service listening on %s\n", addr)
	return server.ListenAndServe()
}
