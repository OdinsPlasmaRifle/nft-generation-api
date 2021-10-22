package main

import (
	"os"
    "net/http"
    "log"
	"time"
	"fmt"
	"encoding/json"
	"image"
	"image/draw"
	"image/png"

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

	if address == "" || id == "" {
		res.HttpCode = http.StatusBadRequest
		res.Data = map[string]interface{}{"message": "A token address and id must be provided."}
		res.RenderJson(w)
		return
	}

	// Build the image.
	var images []image.Image
	mainImage := image.NewRGBA(image.Rect(0, 0, 20, 20))
	for _, img := range images {
		draw.Draw(mainImage, img.Bounds(), img, image.ZP, draw.Over)
	}

	// Create file output.
	file, err := os.Create(fmt.Sprintf("./%s_%s.png", address, id))
	if err != nil {
		panic("Error creating the file.")
	}

	// Encode the image to png in the file.
	err = png.Encode(file, mainImage)
	if err != nil {
		panic("Error encoding the file.")
	}

	res.HttpCode = http.StatusOK
	res.Data = map[string]interface{}{
		"address": address, "id": id,
	}
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