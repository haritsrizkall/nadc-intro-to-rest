package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/ridwanakf/nadc-intro-to-rest/internal"
)

func initFlags(args *internal.Args) {
	port := flag.Int("port", 3000, "port number for your apps")
	args.Port = *port
}

func initHandler(handler *internal.Handler) error {

	// Initialize SQL DB
	db, err := sql.Open("postgres", "postgres://username:password@host-address:port/database-name?sslmode=disable")
	if err != nil {
		return err
	}

	// Check if db connected
	if err = db.PingContext(context.Background()); err != nil {
		return err
	}

	handler.DB = db

	return nil
}

func initRouter(router *httprouter.Router, handler *internal.Handler) {

	router.GET("/", handler.Index)

	// `httprouter` library uses `ServeHTTP` method for it's 404 pages
	router.NotFound = handler
}

func main() {
	args := new(internal.Args)
	initFlags(args)

	handler := new(internal.Handler)
	if err := initHandler(handler); err != nil {
		log.Println("Failed to init handler", err)
		panic(err)
	}

	router := httprouter.New()
	initRouter(router, handler)

	fmt.Printf("Apps served on :%d\n", args.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.Port), router))
}