package main

import (
	"flag"
)

var apiURL string
var useFreshDatabase bool

func init() {
	flag.StringVar(&apiURL, "api-url", "http://localhost:8080", "faker api url")
	flag.BoolVar(&useFreshDatabase, "fresh", false, "use fresh database")

	flag.Parse()
}

func main() {
	done := make(chan bool)

	// Database
	db = NewDB("cpe.db")
	defer db.Close()

	DoMigrations(useFreshDatabase)

	go func() {
		d := NewDispatcher()
		d.Run()
	}()

	go func() {
		c := NewCollector()
		c.Run()
	}()

	go func() {
		r := NewRecorder()
		r.Run()
	}()

	go func() {
		a := NewAnalyzer()
		a.Run()
	}()

	// Router
	go func() {
		r := NewRouter()
		r.Run()
	}()

	<-done
}
