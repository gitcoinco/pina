package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	assert "github.com/gravityblast/miniassert"
)

func TestIndexHandler(t *testing.T) {
	router := newRouter("./")
	handler, _, _ := router.Lookup("GET", "/")
	w := httptest.NewRecorder()
	handler(w, nil, nil)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, "Hello World", string(body))
}

func TestStaticFiles(t *testing.T) {
	publicPath := t.TempDir()
	ipfsPath := filepath.Join(publicPath, "ipfs")
	err := os.Mkdir(ipfsPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	filePath := filepath.Join(ipfsPath, "test.txt")
	err = os.WriteFile(filePath, []byte("Static file example"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	router := newRouter(publicPath)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ipfs/test.txt", nil)
	router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, "Static file example", string(body))
}
