package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	cid "github.com/ipfs/go-cid"
	"github.com/julienschmidt/httprouter"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"
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

func bytesToCID(b []byte) (string, error) {
	pref := cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1, // default length
	}

	hash, err := pref.Sum(b)
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

func handleUpload(w http.ResponseWriter, r *http.Request, content []byte) error {
	// generate CID
	ipfsHash, err := bytesToCID(content)
	if err != nil {
		return err
	}

	// create file using CID as file name
	f, err := os.Create(filepath.Join(ipfsPath, ipfsHash))
	if err != nil {
		return err
	}
	defer f.Close()

	// write uploadded file content to file
	_, err = f.Write(content)
	if err != nil {
		return err
	}

	// encode response body
	err = json.NewEncoder(w).Encode(&PinJSONResponseBody{
		IpfsHash:  ipfsHash,
		PinSize:   10,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	return err
}

func pinJSONHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body PinJSONRequestBody

	// parse request body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		handleError(w, err)
		return
	}

	// encode JSON file content
	content := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(content).Encode(body.PinataContent)
	if err != nil {
		handleError(w, err)
		return
	}

	err = handleUpload(w, r, content.Bytes())
	if err != nil {
		handleError(w, err)
	}
}

func pinFileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(10 << 20) // max 10MB

	// parse uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		handleError(w, err)
		return
	}
	defer file.Close()

	// read file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		handleError(w, err)
		return
	}

	err = handleUpload(w, r, content)
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
	router.POST("/pinning/pinFileToIPFS", pinFileHandler)
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
