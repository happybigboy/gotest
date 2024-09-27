package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Utility function to generate headers
func getHeaders(token string, contentType string) http.Header {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)
	if contentType != "" {
		headers.Set("Content-Type", contentType)
	}
	return headers
}

func GetAccessToken(baseURL, username, password string) (string, error) {
	// Construct the full URL
	reqURL := fmt.Sprintf("%s/api/admin/token", baseURL)
	fmt.Println("Request URL:", reqURL) // Debugging statement

	// Print out the username and password to verify they are set
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
	// Prepare the data as application/x-www-form-urlencoded
	form := url.Values{}
	form.Set("grant_type", "")
	form.Set("username", username)
	form.Set("password", password)
	form.Set("scope", "")
	form.Set("client_id", "")
	form.Set("client_secret", "")

	// Create a new POST request
	req, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("network error occurred: %w", err)
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check the response status
	if res.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		accessToken, ok := result["access_token"].(string)
		if !ok {
			return "", fmt.Errorf("access token not found or invalid type")
		}
		return accessToken, nil
	} else {
		return "", fmt.Errorf("failed to obtain token: %s", body)
	}
}

func GetUserInfo(token, url, username string) (map[string]interface{}, error) {
	userInfoURL := fmt.Sprintf("%s/api/user/%s", url, username)
	headers := getHeaders(token, "")

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		var userInfo map[string]interface{}
		if err := json.Unmarshal(body, &userInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return userInfo, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, &UserNotFoundError{Message: fmt.Sprintf("User '%s' not found", username)}
	}

	return nil, &APIError{Message: fmt.Sprintf("Could not obtain user info for %s: %s", username, body)}
}
