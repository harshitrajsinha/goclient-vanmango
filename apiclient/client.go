package apiclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type APIClient struct {
	client *http.Client
}

func NewAPIClient() *APIClient {
	return &APIClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     30 * time.Second,
				MaxIdleConnsPerHost: 10,
			},
		},
	}
}

// Function to send GET request
func (api *APIClient) Get(url string) (*map[string]interface{}, error) {
	var err error

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	response, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// get response body into byte stream
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.NewDecoder(strings.NewReader(string(content))).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil

}

// Function to send POST request
func (api *APIClient) Post(url string, request *[]byte, cookieVal string) (*map[string]interface{}, error, int) {
	var err error
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(*request))
	if err != nil {
		return nil, err, 500
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cookieVal)
	response, err := api.client.Do(req)
	if err != nil {
		return nil, err, 500
	}

	defer response.Body.Close() // Close the response body when done

	// Read response statusCode
	statusCode := response.StatusCode

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err, statusCode
	}

	var data map[string]interface{}
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&data)
	if err != nil {
		return nil, err, statusCode
	}

	return &data, nil, statusCode

}

// Function to send PUT request
func (api *APIClient) Put(url string, request *[]byte, cookieVal string) (*map[string]interface{}, error) {
	var err error
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(*request))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cookieVal)
	response, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close() // Close the response body when done

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil

}

// Function to send DELETE request
func (api *APIClient) Delete(url string, cookieVal string) (int, error) {
	var err error

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cookieVal)
	response, err := api.client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer response.Body.Close() // Close the response body when done

	// Read status code
	statusCode := response.StatusCode
	// Read the response body
	_, err = io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return statusCode, nil

}
