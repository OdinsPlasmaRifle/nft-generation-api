package main

import (
    "net/http"
    "log"
	"time"
	"encoding/json"

    "github.com/julienschmidt/httprouter"
)

type Response struct {
	HttpCode int
	Data     interface{}
}

func (r *Response) RenderJson(w http.ResponseWriter) {
	successJson, _ := json.MarshalIndent(r.Data, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.HttpCode)
	w.Write(successJson)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{}
	res.HttpCode = http.StatusOK
	res.Data = map[string]interface{}{"message": "NFT Generation API"}
	res.RenderJson(w)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    address := q.Get("address")
    id := q.Get("id")

	res := Response{}

	if (address == "" || id == "") {
		res.HttpCode = http.StatusBadRequest
		res.Data = map[string]interface{}{"message": "A token address and id must be provided."}
		res.RenderJson(w)
		return
	}

	res.HttpCode = http.StatusOK
	res.Data = map[string]interface{}{"address": address, "id": id}
	res.RenderJson(w)
}

func main() {
    router := httprouter.New()
    router.GET(basicChain("/", indexHandler))
    router.GET(basicChain("/token", tokenHandler))

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1048576,
	}

	log.Print("Starting server on 8080.")
	log.Fatal(server.ListenAndServe())
}