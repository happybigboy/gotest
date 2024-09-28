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

func ResetUsage(token, url, username string) error {
	resetURL := fmt.Sprintf("%s/api/user/%s/reset", url, username)
	headers := getHeaders(token, "")

	req, err := http.NewRequest("POST", resetURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil // Reset was successful
	}

	return &APIError{Message: fmt.Sprintf("Could not reset usage for %s: %s", username, body)}
}

func RevokeSubscription(token, url, username string) (string, error) {
    revokeURL := fmt.Sprintf("%s/api/user/%s/revoke_sub", url, username)
    headers := getHeaders(token, "application/json")
    data := map[string]interface{}{} // Empty JSON object

    reqBody, err := json.Marshal(data)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request body: %w", err)
    }

    req, err := http.NewRequest("POST", revokeURL, bytes.NewBuffer(reqBody))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        // Assuming the subscription link is returned in the response
        var response map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
            return "", fmt.Errorf("failed to decode response body: %w", err)
        }
        if link, ok := response["subscription_url"].(string); ok {
            return link, nil // Return the subscription link
        }
        return "", fmt.Errorf("no subscription link found in the response")
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response body: %w", err)
    }

    return "", &APIError{Message: fmt.Sprintf("Could not revoke subscription for %s: %s", username, body)}
}


func GetInbounds(token, url string) (interface{}, error) {
	inboundsURL := fmt.Sprintf("%s/api/inbounds", url)
	headers := getHeaders(token, "")

	req, err := http.NewRequest("GET", inboundsURL, nil)
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
		var result interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("could not obtain inbounds: %s", string(body))
	}
}

func CreateUser(token, url, username string, expire int64, limit int64, note string) (map[string]interface{}, error) {
	headers := getHeaders(token, "application/json")
	userURL := fmt.Sprintf("%s/api/user/", url)

	inboundsData, err := GetInbounds(token, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbounds: %w", err)
	}

	proxies := make(map[string]map[string]interface{})
	inbounds := make(map[string][]string)

	// Process inbound data
	inboundsMap, ok := inboundsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("inbounds data is not in expected format")
	}

	for _, inboundList := range inboundsMap {
		// Ensure inboundList is a slice of interfaces
		if inboundArray, ok := inboundList.([]interface{}); ok {
			for _, inbound := range inboundArray {
				// Ensure inbound is a map
				if inboundMap, ok := inbound.(map[string]interface{}); ok {
					protocol := inboundMap["protocol"].(string)
					tag := inboundMap["tag"].(string)
	
					if protocol != "" {
						// Initialize the proxies map if it doesn't exist
						if _, exists := proxies[protocol]; !exists {
							proxies[protocol] = make(map[string]interface{})
						}
						// Initialize the inbounds slice if it doesn't exist
						if _, exists := inbounds[protocol]; !exists {
							inbounds[protocol] = []string{}
						}
						// Append the tag if it's not empty
						if tag != "" {
							inbounds[protocol] = append(inbounds[protocol], tag)
						}
					}
				}
			}
		}
	}

	data := map[string]interface{}{
		"username":                       username,
		"proxies":                        proxies,
		"inbounds":                       inbounds,
		"expire":                         expire,
		"data_limit":                     limit,
		"data_limit_reset_strategy":      "no_reset",
		"status":                         "active",
		"note":                           note,
		"on_hold_timeout":                nil,
		"on_hold_expire_duration":        0,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequest("POST", userURL, bytes.NewBuffer(dataBytes))
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

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}
		return result, nil
	}

	var errorMessage string
	if err := json.Unmarshal(body, &errorMessage); err != nil {
		errorMessage = string(body) // Fallback to raw response if JSON unmarshal fails
	}

	if len(errorMessage) > 500 {
		errorMessage = errorMessage[:500] + "... [truncated]"
	}
	return nil, fmt.Errorf("could not create user %s: %s", username, errorMessage)
}
func ModifyUser(token, url, username string, expire int64, limit int64, note string, status string) (map[string]interface{}, error) {
	headers := getHeaders(token, "application/json")
	userURL := fmt.Sprintf("%s/api/user/%s", url, username)

	inboundsData, err := GetInbounds(token, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbounds: %w", err)
	}

	proxies := make(map[string]map[string]interface{})
	inbounds := make(map[string][]string)

	// Process inbound data
	inboundsMap, ok := inboundsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("inbounds data is not in expected format")
	}

	for _, inboundList := range inboundsMap {
		if inboundArray, ok := inboundList.([]interface{}); ok {
			for _, inbound := range inboundArray {
				inboundMap, ok := inbound.(map[string]interface{})
				if !ok {
					continue
				}
				protocol, _ := inboundMap["protocol"].(string)
				tag, _ := inboundMap["tag"].(string)

				if protocol != "" {
					if _, exists := proxies[protocol]; !exists {
						proxies[protocol] = make(map[string]interface{})
					}
					if _, exists := inbounds[protocol]; !exists {
						inbounds[protocol] = []string{}
					}
					if tag != "" {
						inbounds[protocol] = append(inbounds[protocol], tag)
					}
				}
			}
		}
	}

	data := map[string]interface{}{
		"proxies":                       proxies,
		"inbounds":                      inbounds,
		"expire":                        expire,
		"data_limit":                    limit,
		"data_limit_reset_strategy":     "no_reset",
		"status":                        status,
		"note":                          note,
		"on_hold_timeout":               nil,
		"on_hold_expire_duration":       0,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequest("PUT", userURL, bytes.NewBuffer(dataBytes))
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

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}
		return result, nil
	}

	var errorMessage string
	if err := json.Unmarshal(body, &errorMessage); err != nil {
		errorMessage = string(body) // Fallback to raw response if JSON unmarshal fails
	}

	if len(errorMessage) > 500 {
		errorMessage = errorMessage[:500] + "... [truncated]"
	}
	return nil, fmt.Errorf("could not modify user %s: %s", username, errorMessage)
}

func GetUsers(token, baseURL string, offset, limit int, sort string) (map[string]interface{}, error) {
	client := &http.Client{}
	headers := getHeaders(token, "")
	usersURL, err := buildURL(baseURL, offset, limit, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to build request URL: %w", err)
	}

	req, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return nil, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := checkResponse(resp, body); err != nil {
		return nil, err
	}

	return parseResponse(body)
}

func GetAllUsers(token, baseURL string, offset, limit int, sort string) (int, error) {
	client := &http.Client{}
	headers := getHeaders(token, "")
	usersURL, err := buildURL(baseURL, offset, limit, sort)
	if err != nil {
		return -1, fmt.Errorf("failed to build request URL: %w", err)
	}

	req, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return 0, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := checkResponse(resp, body); err != nil {
		return 0, err
	}

	var responseMap struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return responseMap.Total, nil
}
func buildURL(baseURL string, offset, limit int, sort string) (string, error) {
	params := url.Values{}
	params.Add("offset", fmt.Sprintf("%d", offset))
	params.Add("limit", fmt.Sprintf("%d", limit))
	if sort != "" {
		params.Add("sort", sort)
	}
	return fmt.Sprintf("%s/api/users?%s", baseURL, params.Encode()), nil
}

func checkResponse(resp *http.Response, body []byte) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not obtain users: %s", string(body))
	}
	return nil
}

func parseResponse(body []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	usersList, ok := result["users"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'users' field to be an array but got: %v", result["users"])
	}
	total := result["total"]
	if total == nil {
		total = len(usersList) // Fallback to count if total is not provided
	}

	return map[string]interface{}{
		"total": total,
		"users": usersList,
	}, nil
}
