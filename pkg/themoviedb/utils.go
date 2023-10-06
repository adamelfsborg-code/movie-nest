package themoviedb

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	baseUrl = "https://api.themoviedb.org"
)

type MovieDBOptions struct {
	AuthToken  string `json:"authtoken"`
	ApiKey     string `json:"apikey"`
	ApiVersion uint8  `json:"version"`
}

type LoggingTransport struct {
	Transport http.RoundTripper
}

func NewMovieDBOptions(authToken string, apiKey string, apiVersion ...uint8) *MovieDBOptions {
	mo := &MovieDBOptions{
		AuthToken:  authToken,
		ApiKey:     apiKey,
		ApiVersion: 3,
	}

	if len(apiVersion) > 0 {
		mo.ApiVersion = apiVersion[0]
	}

	return mo
}

func (m *MovieDBOptions) setHeaders(request *http.Request) error {
	authKey := ""

	if m.AuthToken != "" {
		authKey = m.AuthToken
	}

	if m.ApiKey != "" {
		authKey = m.ApiKey
	}

	if authKey == "" {
		return fmt.Errorf("AuhToken or ApiKey not set")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("accept", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", authKey))

	return nil
}

func (m *MovieDBOptions) Get(path string) ([]byte, error) {
	client := &http.Client{
		Transport: &LoggingTransport{
			Transport: http.DefaultTransport,
		},
	}

	url := fmt.Sprintf("%v/%v/%v", baseUrl, m.ApiVersion, path)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	err = m.setHeaders(request)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return body, nil
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	log.Printf("Outgoing HTTP request: %s %s", req.Method, req.URL.String())

	logRequestHeaders(req)

	logRequestBody(req)

	resp, err := t.Transport.RoundTrip(req)

	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return nil, err
	}

	elapsed := time.Since(start)
	log.Printf("HTTP response status: %s", resp.Status)

	logResponseHeaders(resp)

	log.Printf("HTTP request took: %v", elapsed)

	return resp, nil
}

func logRequestHeaders(req *http.Request) {
	log.Println("Request Headers:")
	for key, values := range req.Header {
		for _, value := range values {
			log.Printf("%s: %s", key, value)
		}
	}
}

func logRequestBody(req *http.Request) {
	if req.Body == nil {
		log.Println("Request Body: <empty>")
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return
	}

	log.Printf("Request Body: %s", body)

	req.Body = io.NopCloser(bytes.NewReader(body))
}

func logResponseHeaders(resp *http.Response) {
	log.Println("Response Headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			log.Printf("%s: %s", key, value)
		}
	}
}
