package main

import (
	"io"
	"net/http/httptest"
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
