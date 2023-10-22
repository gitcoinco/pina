package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

type PinJSONRequestBody struct {
	PinataContent interface{} `json:"pinataContent"`
}

type PinJSONResponseBody struct {
	IpfsHash  string `json:"ipfsHash"`
	PinSize   int    `json:"pinSize"`
	Timestamp string `json:"timestamp"`
}

var (
	publicPath string
	ipfsPath   string
	port       int
)

func init() {
	flag.IntVar(&port, "port", 0, "http server port")
	flag.StringVar(&publicPath, "public", "", "public path")
}
func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "server error1")
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World")
}

func pinJSONHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body PinJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		handleError(w, err)
		return
	}

	fileName := "test.json"
	f, err := os.Create(filepath.Join(ipfsPath, fileName))
	if err != nil {
		handleError(w, err)
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(body.PinataContent)
	if err != nil {
		handleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&PinJSONResponseBody{
		IpfsHash:  "foo",
		PinSize:   10,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	if err != nil {
		handleError(w, err)
	}
}

func newRouter(publicPath string) (*httprouter.Router, error) {
	ipfsPath = filepath.Join(publicPath, "ipfs")
	err := os.MkdirAll(ipfsPath, 0755)
	if err != nil {
		return nil, err
	}

	router := httprouter.New()
	router.GET("/", indexHandler)
	router.POST("/pinning/pinJSONToIPFS", pinJSONHandler)
	router.NotFound = http.FileServer(http.Dir(publicPath))

	return router, nil
}

func main() {
	flag.Parse()
	if port == 0 || publicPath == "" {
		fmt.Println("port and public path flags are mandatory")
		flag.Usage()
		os.Exit(1)
	}

	router, err := newRouter(publicPath)
	if err != nil {
		log.Fatal(err)
	}
	binding := fmt.Sprintf(":%d", port)
	fmt.Printf("listening: %s\n", binding)
	fmt.Printf("public path: %s\n", publicPath)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
