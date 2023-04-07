package simplehttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// SimpleHttpResponse - a simple type that contains http responses
type SimpleHttpResponse struct {
	StatusCode int
	Body       []byte
}

var simpleClient http.Client

// init - Initilize client on startup
func init() {
	simpleClient = http.Client{Timeout: (20 * time.Second)}
}

// SimpleGet - A simple get request
func SimpleGet(uri string, headers map[string]string) (*SimpleHttpResponse, error) {
	// Parse the URL to check if it contains query parameters
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	// If the URL contains query parameters, escape them
	if parsedURL.RawQuery != "" {
		parsedURL.RawQuery = url.PathEscape(parsedURL.RawQuery)
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %s", err.Error())
	}

	//Add headers from map
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request and get the response
	resp, err := simpleClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Process response into type
	var response SimpleHttpResponse
	response.StatusCode = resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response.Body = body

	return &response, nil
}
