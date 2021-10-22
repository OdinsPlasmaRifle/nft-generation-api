package main

import (
	"os"
    "net/http"
    "log"
	"time"
	"fmt"
	"io/ioutil"
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
		res := Response{
			HttpCode: http.StatusNotFound,
			Data:     map[string]interface{}{"message": http.StatusText(http.StatusNotFound)},
		}
		res.RenderJson(w)
		return
	}

	// Get list of assets for the image.
	// TODO : needs to figure these out based on the smart contract.
	paths := []string{
		"./assets/backgrounds/Bg-blue.png",
		"./assets/skin/Base-F-1.png",
		"./assets/outfits/Outfit1.png",
		"./assets/hair/Hair-blonde.png",
		"./assets/eyes/Eyes-blue.png",
		"./assets/lips/Lips-orange.png",
		"./assets/accessory/Acc-earring-gold.png",
	}
	var assets []image.Image
	for _, path := range paths {
		f, err := os.Open(path) 
		if err != nil {
			panic("Error opening an asset file.")
		}
		defer f.Close()

		asset, err := png.Decode(f)
		if err != nil {
			panic("Error decoding an asset file.")
		}

		assets = append(assets, asset)
	}

	// Build the image using the assets.
	mainImage := image.NewRGBA(image.Rect(0, 0, 1336, 1336))
	for _, asset := range assets {
		draw.Draw(mainImage, asset.Bounds(), asset, image.ZP, draw.Over)
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

func ImageHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    address := q.Get("address")
    id := q.Get("id")

	if address == "" || id == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "appplication/octet-stream")
		return
	}

	fileBytes, err := ioutil.ReadFile(fmt.Sprintf("./%s_%s.png", address, id))
	if err != nil {
		panic("Error fetching the file.")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "appplication/octet-stream")
	w.Write(fileBytes)
}

func main() {
    router := httprouter.New()
    router.GET(basicChain("/", indexHandler))
    router.GET(basicChain("/token", tokenHandler))
	router.GET(basicChain("/image", ImageHandler))

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