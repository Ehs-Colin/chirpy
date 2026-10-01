package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	myHandler := http.NewServeMux()
	myHandler.Handle("/", http.FileServer(http.Dir(".")))
	s := &http.Server{
		Addr:           ":8080",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
