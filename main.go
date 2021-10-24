package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
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

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	res := Response{
		HttpCode: http.StatusNotFound,
		Data:     map[string]interface{}{"message": http.StatusText(http.StatusNotFound)},
	}
	res.RenderJson(w)
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
	id, err := strconv.Atoi(q.Get("id"))
	res := Response{}

	// Throw an error is the address or ID are not set.
	if err != nil || address == "" {
		notFoundHandler(w, r)
		return
	}

	// Fetch the contract using the address.
	contract, err := NewERC1155(address)
	if err != nil {
		panic(err)
	}
	// Collect the attribute using the id.
	attributes, err := contract.GetTokenAttributes(id)
	if err != nil {
		panic(err)
	}

	// Get the assets for the attributes.
	assets, _ := GetAssets(attributes)
	// Generate an image using the assets.
	image, _ := CreateImage(assets)
	// Create a file for the image
	file, _ := CreatePNGFile(image, fmt.Sprintf("%s_%d", address, id))
	// Create a file path for the image.
	filePath := strings.Join(
		[]string{"https://", r.Host, "/images/", path.Base(file.Name())}, "",
	)

	res.HttpCode = http.StatusOK
	res.Data = map[string]interface{}{
		"image": filePath, "attributes": attributes,
	}
	res.RenderJson(w)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	ps := r.Context().Value("params").(httprouter.Params)
	fileBytes, err := ioutil.ReadFile(fmt.Sprintf("./var/images/%s", ps.ByName("name")))
	if err != nil {
		// Set a non JSON not found error.
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "appplication/octet-stream")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "appplication/octet-stream")
	w.Write(fileBytes)
}

func main() {
	chain := alice.New(
		loggerHandler, recoveryHandler, corsHandler,
	)
	router := httprouter.New()
	router.GET("/", wrapper(chain.ThenFunc(indexHandler)))
	router.GET("/token", wrapper(chain.ThenFunc(tokenHandler)))
	router.GET("/images/:name", wrapper(chain.ThenFunc(imageHandler)))

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
