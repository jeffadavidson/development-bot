package simplehttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HttpGet_ProcessesBody(t *testing.T) {
	// Create a test server
	expectedBody := `{"name":"John", "age":30}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		// Set the response body
		fmt.Fprintf(w, `{"name":"John", "age":30}`)
	}))
	defer ts.Close()

	response, err := SimpleGet(ts.URL, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, expectedBody, string(response.Body))
}

func Test_HttpGet_ProcessesEmptyBody(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		// Set the response body
		fmt.Fprintf(w, "")
	}))
	defer ts.Close()

	response, err := SimpleGet(ts.URL, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "", string(response.Body))
}

func Test_HttpGet_404IsNotTreatedAsError(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer ts.Close()

	response, err := SimpleGet(ts.URL, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
