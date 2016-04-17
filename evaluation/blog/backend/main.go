package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	s = newStore(*dbFile)

	http.HandleFunc("/api/posts", apiPosts)

	http.Handle("/", http.FileServer(http.Dir("frontend")))

	log.Fatal(http.ListenAndServe(*httpSpec, nil))
}
