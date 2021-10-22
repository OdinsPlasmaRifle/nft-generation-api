package main

import (
    "fmt"
    "net/http"
    "log"

    "github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    "0m5s",
		WriteTimeout:   "0m5s",
		MaxHeaderBytes: 1048576,
	}

	log.Print("Starting server on 8080.")

	log.Fatal(server.ListenAndServe())
}