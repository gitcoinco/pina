package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	assert "github.com/gravityblast/miniassert"
	"github.com/julienschmidt/httprouter"
)

func newTestRouter(t *testing.T) (*httprouter.Router, string) {
	publicPath := t.TempDir()
	router, err := newRouter(publicPath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return router, filepath.Join(publicPath, "ipfs")
}

func TestIndexHandler(t *testing.T) {
	router, _ := newTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, "Hello World", string(body))
}

func TestStaticFiles(t *testing.T) {
	router, ipfsPath := newTestRouter(t)
	filePath := filepath.Join(ipfsPath, "test.txt")
	err := os.WriteFile(filePath, []byte("Static file example"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ipfs/test.txt", nil)
	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, "Static file example", string(body))
}

func TestPinJSONHandler(t *testing.T) {
	router, ipfsPath := newTestRouter(t)

	w := httptest.NewRecorder()
	reqBody := bytes.NewBufferString(`{"pinataContent": {"foo": {"bar": "baz"}}}`)
	req := httptest.NewRequest(http.MethodPost, "/pinning/pinJSONToIPFS", reqBody)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", AUTH_TOKEN))
	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

	var responseBody PinJSONResponseBody
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		log.Fatal(err)
	}

	cid := "bafkreihktyturq4bzrikdjylvvjbrgh5rfigzlydmjoyri3ip6fjbcqddu"
	assert.Equal(t, cid, responseBody.IpfsHash)
	assert.Equal(t, 10, responseBody.PinSize)
	assert.NotEqual(t, "", responseBody.Timestamp)

	fileName := filepath.Join(ipfsPath, cid)
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, `{"foo":{"bar":"baz"}}`, strings.TrimSpace(string(content)))
}

func TestPinFileHandler(t *testing.T) {
	router, _ := newTestRouter(t)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// writer.WriteField("file", "Uploaded content")

	// Create the file part with an invalid file extension
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "uploaded.txt"))
	h.Set("Content-Type", "text/plain")
	part, _ := writer.CreatePart(h)
	part.Write([]byte("Uploaded content"))
	writer.Close()
	// Prepare and send the HTTP request

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/pinning/pinFileToIPFS", body)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", AUTH_TOKEN))
	req.Header.Add("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

	var responseBody PinJSONResponseBody
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		log.Fatal(err)
	}

	cid := "bafkreifqrpvngn6k2q6qahrm2oawbtrsoucxi7xtxehydodbzgnky6eiem"
	assert.Equal(t, cid, responseBody.IpfsHash)
	assert.Equal(t, 10, responseBody.PinSize)
	assert.NotEqual(t, "", responseBody.Timestamp)

	fileName := filepath.Join(ipfsPath, cid)
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, "Uploaded content", strings.TrimSpace(string(content)))
}
