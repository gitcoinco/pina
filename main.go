package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var (
	publicPath string
	port       int
)

func init() {
	flag.IntVar(&port, "port", 0, "http server port")
	flag.StringVar(&publicPath, "public", "", "public path")
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World")
}

func newRouter(publicPath string) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.NotFound = http.FileServer(http.Dir(publicPath))

	return router
}

func main() {
	flag.Parse()
	if port == 0 || publicPath == "" {
		fmt.Println("port and public path flags are mandatory")
		flag.Usage()
		os.Exit(1)
	}

	router := newRouter(publicPath)
	binding := fmt.Sprintf(":%d", port)
	fmt.Printf("listening: %s\n", binding)
	fmt.Printf("public path: %s\n", publicPath)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
