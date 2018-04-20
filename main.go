package main

import (
	"flag"
)

var doSeed bool

func init() {
	flag.Parse()
}

func main() {
	// Database
	db = NewDB("cpe.db")
	defer db.Close()

	// Router
	r := NewRouter()
	r.Run()
}
