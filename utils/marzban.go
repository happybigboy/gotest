package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Makes a GET request to https://localhost:3000/index.php and returns the JSON response as a string
func MakeRequest() (string, error) {
	// Make a GET request to the local server
	resp, err := http.Get("https://localhost:3000/index.php")
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert the JSON response to a string
	jsonString, err := json.MarshalIndent(jsonResponse, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonString), nil
}
