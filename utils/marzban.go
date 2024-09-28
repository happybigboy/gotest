package utils

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "io/ioutil"
	// "net/http"
	// "net/url"
)

// Custom error types
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

type NetworkError struct {
	Message string
}

func (e *NetworkError) Error() string {
	return e.Message
}

type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

type UserNotFoundError struct {
	Message string
}

func (e *UserNotFoundError) Error() string {
	return e.Message
}



// // GetAccessToken retrieves the access token using the provided credentials.
// func GetAccessToken(username, password, apiUrl string) (string, error) {
// 	loginURL := fmt.Sprintf("%s/api/admin/token/", apiUrl)

// 	// Prepare the data in application/x-www-form-urlencoded format
// 	data := url.Values{}
// 	data.Set("username", username)
// 	data.Set("password", password)

// 	resp, err := http.PostForm(loginURL, data)
// 	if err != nil {
// 		return "", &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	if resp.StatusCode == http.StatusOK {
// 		var result map[string]interface{}
// 		if err := json.Unmarshal(body, &result); err != nil {
// 			return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
// 		}

// 		accessToken, ok := result["access_token"].(string)
// 		if !ok {
// 			return "", fmt.Errorf("access token not found or invalid type")
// 		}
// 		return accessToken, nil
// 	} else if resp.StatusCode == http.StatusUnauthorized {
// 		return "", &AuthError{Message: "Invalid username or password"}
// 	}

// 	return "", &APIError{Message: fmt.Sprintf("Could not obtain token: %s", body)}
// }

// // GetUserInfo retrieves user information for the specified username.
// func getUserInfo(token, url, username string) (map[string]interface{}, error) {
// 	userInfoURL := fmt.Sprintf("%s/api/user/%s", url, username)
// 	headers := getHeaders(token, "")

// 	req, err := http.NewRequest("GET", userInfoURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header = headers

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode == 200 {
// 		var userInfo map[string]interface{}
// 		if err := json.Unmarshal(body, &userInfo); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
// 		}
// 		return userInfo, nil
// 	} else if resp.StatusCode == 404 {
// 		return nil, &UserNotFoundError{Message: fmt.Sprintf("User '%s' not found", username)}
// 	}
// 	return nil, &APIError{Message: fmt.Sprintf("Could not obtain user info for %s: %s", username, body)}
// }

// ResetUsage resets the usage for the specified user.
// func resetUsage(token, url, username string) (bool, error) {
// 	resetURL := fmt.Sprintf("%s/api/user/%s/reset", url, username)
// 	headers := getHeaders(token, "")

// 	req, err := http.NewRequest("POST", resetURL, nil)
// 	if err != nil {
// 		return false, err
// 	}
// 	req.Header = headers

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return false, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode == 200 {
// 		return true, nil
// 	}
// 	return false, fmt.Errorf("Could not reset usage for %s: %s", username, body)
// }

// // RevokeSubscription revokes the subscription for the specified user.
// func revokeSubscription(token, url, username string) (bool, error) {
// 	revokeURL := fmt.Sprintf("%s/api/user/%s/revoke_sub", url, username)
// 	headers := getHeaders(token, "application/json")

// 	data := map[string]interface{}{}
// 	jsonData, _ := json.Marshal(data)

// 	req, err := http.NewRequest("POST", revokeURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return false, err
// 	}
// 	req.Header = headers

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return false, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode == 200 {
// 		return true, nil
// 	}
// 	return false, fmt.Errorf("Could not revoke subscription for %s: %s", username, body)
// }

// // GetInbounds retrieves inbound information.
// func getInbounds(token, url string) (map[string]interface{}, error) {
// 	inboundsURL := fmt.Sprintf("%s/api/inbounds", url)
// 	headers := getHeaders(token, "")

// 	req, err := http.NewRequest("GET", inboundsURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header = headers

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode == 200 {
// 		var result map[string]interface{}
// 		if err := json.Unmarshal(body, &result); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
// 		}
// 		return result, nil
// 	}
// 	return nil, fmt.Errorf("Could not obtain inbounds: %s", body)
// }

// func createUser(token, url, username, expire string, limit int, note string) (map[string]interface{}, error) {
// 	headers := getHeaders(token, "application/json")
// 	userURL := fmt.Sprintf("%s/api/user/", url)

// 	inboundsData, err := getInbounds(token, url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	proxies := make(map[string]interface{})
// 	inbounds := make(map[string]interface{})

// 	// Extracting proxies and inbounds from the inbound data
// 	for _, inboundList := range inboundsData {
// 		for _, inbound := range inboundList.([]interface{}) {
// 			protocol := inbound.(map[string]interface{})["protocol"]
// 			tag := inbound.(map[string]interface{})["tag"]
// 			if protocol != nil {
// 				if _, exists := proxies[protocol.(string)]; !exists {
// 					proxies[protocol.(string)] = struct{}{}
// 				}
// 				if _, exists := inbounds[protocol.(string)]; !exists {
// 					inbounds[protocol.(string)] = []string{}
// 				}
// 				if tag != nil {
// 					inbounds[protocol.(string)] = append(inbounds[protocol.(string)].([]string), tag.(string))
// 				}
// 			}
// 		}
// 	}

// 	data := map[string]interface{}{
// 		"username":                  username,
// 		"proxies":                   proxies,
// 		"inbounds":                  inbounds,
// 		"expire":                    expire,
// 		"data_limit":                limit,
// 		"data_limit_reset_strategy": "no_reset",
// 		"status":                    "active",
// 		"note":                      note,
// 		"on_hold_timeout":           nil,
// 		"on_hold_expire_duration":   0,
// 	}

// 	// Marshal the data into JSON
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
// 	}

// 	// Create a new POST request with headers
// 	req, err := http.NewRequest("POST", userURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header = headers // Setting the headers

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode == 200 || resp.StatusCode == 201 {
// 		var result map[string]interface{}
// 		if err := json.Unmarshal(body, &result); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
// 		}
// 		return result, nil
// 	}

// 	return nil, fmt.Errorf("Could not create user %s: %s", username, body)
// }

// // ModifyUser modifies an existing user's details.
// func modifyUser(token, url, username, expire string, limit int, note string, status string) (map[string]interface{}, error) {
// 	headers := getHeaders(token, "application/json")
// 	userURL := fmt.Sprintf("%s/api/user/%s", url, username)

// 	inboundsData, err := getInbounds(token, url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	proxies := make(map[string]interface{})
// 	inbounds := make(map[string]interface{})

// 	for _, inboundList := range inboundsData {
// 		for _, inbound := range inboundList.([]interface{}) {
// 			protocol := inbound.(map[string]interface{})["protocol"]
// 			tag := inbound.(map[string]interface{})["tag"]
// 			if protocol != nil {
// 				if _, exists := proxies[protocol.(string)]; !exists {
// 					proxies[protocol.(string)] = struct{}{}
// 				}
// 				if _, exists := inbounds[protocol.(string)]; !exists {
// 					inbounds[protocol.(string)] = []string{}
// 				}
// 				if tag != nil {
// 					inbounds[protocol.(string)] = append(inbounds[protocol.(string)].([]string), tag.(string))
// 				}
// 			}
// 		}
// 	}

// 	data := map[string]interface{}{
// 		"proxies":                   proxies,
// 		"inbounds":                  inbounds,
// 		"expire":                    expire,
// 		"data_limit":                limit,
// 		"data_limit_reset_strategy": "no_reset",
// 		"status":                    status,
// 		"note":                      note,
// 		"on_hold_timeout":           nil,
// 		"on_hold_expire_duration":   0,
// 	}

// 	jsonData, _ := json.Marshal(data)
// 	req, err := http.NewRequest("PUT", userURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header = headers

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode == 200 || resp.StatusCode == 201 {
// 		var result map[string]interface{}
// 		if err := json.Unmarshal(body, &result); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
// 		}
// 		return result, nil
// 	}

// 	return nil, fmt.Errorf("Could not modify user %s: %s", username, body)
// }

// GetUsers retrieves a list of users based on the specified parameters.
// func GetUsers(token, url string, offset int, limit int, sort string) (map[string]interface{}, error) {
// 	headers := getHeaders(token, "")
// 	usersURL := fmt.Sprintf("%s/api/users", url)

// 	req, err := http.NewRequest("GET", usersURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header = headers

// 	query := req.URL.Query()
// 	query.Add("offset", fmt.Sprintf("%d", offset))
// 	query.Add("limit", fmt.Sprintf("%d", limit))
// 	query.Add("sort", sort)
// 	req.URL.RawQuery = query.Encode()

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, &NetworkError{Message: fmt.Sprintf("Network error occurred: %s", err)}
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	if resp.StatusCode == 200 {
// 		var result map[string]interface{}
// 		if err := json.Unmarshal(body, &result); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
// 		}
// 		if users, ok := result["users"].([]interface{}); ok {
// 			totalUsers := len(users)
// 			return map[string]interface{}{"total": totalUsers, "users": users}, nil
// 		}
// 		return nil, fmt.Errorf("unexpected response format: %s", body)
// 	}
// 	return nil, fmt.Errorf("Could not obtain users: %s", body)
// }
