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
	"strconv"

	"github.com/julienschmidt/httprouter"
	
	web3 "github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"github.com/umbracle/go-web3/abi"
	"github.com/umbracle/go-web3/contract"
)

type Response struct {
	HttpCode int
	Data     interface{}
}

type ERC1155 struct {
	c *contract.Contract
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
	check_id, err := strconv.Atoi(id)

	res := Response{}

	if address == "" || id == "" {
		res := Response{
			HttpCode: http.StatusNotFound,
			Data:     map[string]interface{}{"message": http.StatusText(http.StatusNotFound)},
		}
		res.RenderJson(w)
		return
	}
	
	// Get contract details
	var erc1155ABI = `[ { "inputs": [], "stateMutability": "nonpayable", "type": "constructor" }, { "anonymous": false, "inputs": [ { "indexed": true, "internalType": "address", "name": "owner", "type": "address" }, { "indexed": true, "internalType": "address", "name": "approved", "type": "address" }, { "indexed": true, "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "Approval", "type": "event" }, { "anonymous": false, "inputs": [ { "indexed": true, "internalType": "address", "name": "owner", "type": "address" }, { "indexed": true, "internalType": "address", "name": "operator", "type": "address" }, { "indexed": false, "internalType": "bool", "name": "approved", "type": "bool" } ], "name": "ApprovalForAll", "type": "event" }, { "anonymous": false, "inputs": [ { "indexed": true, "internalType": "address", "name": "from", "type": "address" }, { "indexed": true, "internalType": "address", "name": "to", "type": "address" }, { "indexed": true, "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "Transfer", "type": "event" }, { "inputs": [ { "internalType": "address", "name": "to", "type": "address" }, { "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "approve", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [ { "internalType": "address", "name": "owner", "type": "address" } ], "name": "balanceOf", "outputs": [ { "internalType": "uint256", "name": "", "type": "uint256" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "getApproved", "outputs": [ { "internalType": "address", "name": "", "type": "address" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "getAttributes", "outputs": [ { "internalType": "string[]", "name": "", "type": "string[]" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "address", "name": "owner", "type": "address" }, { "internalType": "address", "name": "operator", "type": "address" } ], "name": "isApprovedForAll", "outputs": [ { "internalType": "bool", "name": "", "type": "bool" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "address", "name": "to", "type": "address" }, { "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "mint", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [], "name": "name", "outputs": [ { "internalType": "string", "name": "", "type": "string" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "ownerOf", "outputs": [ { "internalType": "address", "name": "", "type": "address" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "address", "name": "from", "type": "address" }, { "internalType": "address", "name": "to", "type": "address" }, { "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "safeTransferFrom", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [ { "internalType": "address", "name": "from", "type": "address" }, { "internalType": "address", "name": "to", "type": "address" }, { "internalType": "uint256", "name": "tokenId", "type": "uint256" }, { "internalType": "bytes", "name": "_data", "type": "bytes" } ], "name": "safeTransferFrom", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [ { "internalType": "address", "name": "operator", "type": "address" }, { "internalType": "bool", "name": "approved", "type": "bool" } ], "name": "setApprovalForAll", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [ { "internalType": "bytes4", "name": "interfaceId", "type": "bytes4" } ], "name": "supportsInterface", "outputs": [ { "internalType": "bool", "name": "", "type": "bool" } ], "stateMutability": "view", "type": "function" }, { "inputs": [], "name": "symbol", "outputs": [ { "internalType": "string", "name": "", "type": "string" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "uint256", "name": "index", "type": "uint256" } ], "name": "tokenByIndex", "outputs": [ { "internalType": "uint256", "name": "", "type": "uint256" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "address", "name": "owner", "type": "address" }, { "internalType": "uint256", "name": "index", "type": "uint256" } ], "name": "tokenOfOwnerByIndex", "outputs": [ { "internalType": "uint256", "name": "", "type": "uint256" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "tokenURI", "outputs": [ { "internalType": "string", "name": "", "type": "string" } ], "stateMutability": "view", "type": "function" }, { "inputs": [], "name": "totalSupply", "outputs": [ { "internalType": "uint256", "name": "", "type": "uint256" } ], "stateMutability": "view", "type": "function" }, { "inputs": [ { "internalType": "address", "name": "from", "type": "address" }, { "internalType": "address", "name": "to", "type": "address" }, { "internalType": "uint256", "name": "tokenId", "type": "uint256" } ], "name": "transferFrom", "outputs": [], "stateMutability": "nonpayable", "type": "function" } ]`
	abiERC1155, err := abi.NewABI(erc1155ABI)
	// Initiate the RPC web3 client
	client, err := jsonrpc.NewClient("PUT_YOUR_INFURA_URL_HERE")
	if err != nil {
		panic(err)
	}
	// Get the latest block
	number, err := client.Eth().BlockNumber()
	if err != nil {
		panic(err)
	}
	// Call the contract to get the attributes
	var contract = ERC1155{c: contract.NewContract(web3.HexToAddress(address), abiERC1155, client)}
	contract_output, err := contract.c.Call("getAttributes", web3.EncodeBlock(web3.BlockNumber(number)), check_id)

	// fmt.Println(attributes["0"][0])
	// s := strings.Split(attributes["0"].(string), " ")
	attributes := contract_output["0"].([]string)
	// fmt.Println(retval0[0])

	// Get list of assets for the image.
	// TODO : needs to figure these out based on the smart contract.
	// paths := []string{
	// 	"./assets/backgrounds/Bg-blue.png",
	// 	"./assets/skin/Base-F-3.png",
	// 	"./assets/outfits/Outfit2.png",
	// 	"./assets/hair/Hair-blonde.png",
	// 	"./assets/eyes/Eyes-blue.png",
	// 	"./assets/lips/Lips-orange.png",
	// 	"./assets/accessory/Acc-earring-gold.png",
	// }
	paths := []string{
		fmt.Sprintf("./assets/backgrounds/Bg-%s.png", attributes[1]),
		fmt.Sprintf("./assets/skin/Base-%s.png", attributes[6]),
		fmt.Sprintf("./assets/outfits/Outfit%s.png", attributes[5]),
		fmt.Sprintf("./assets/hair/Hair-%s.png", attributes[3]),
		fmt.Sprintf("./assets/eyes/Eyes-%s.png", attributes[2]),
		fmt.Sprintf("./assets/lips/Lips-%s.png", attributes[4]),
		fmt.Sprintf("./assets/accessory/Acc-%s.png", attributes[0]),
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