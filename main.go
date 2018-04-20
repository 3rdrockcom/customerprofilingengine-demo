package main

import (
	"flag"
	"runtime"
)

var apiURL string

func init() {
	flag.StringVar(&apiURL, "api-url", "http://localhost:8080", "faker api url")

	flag.Parse()
}

func main() {
	// Database
	db = NewDB("cpe.db")
	defer db.Close()

	go func() {
		d := NewDispatcher()
		d.Run()
	}()

	// Router
	go func() {
		r := NewRouter()
		r.Run()
	}()

	// Keep the connection alive
	runtime.Goexit()
}
