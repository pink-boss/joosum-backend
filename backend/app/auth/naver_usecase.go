package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NaverUsecae struct {
}

func (n *NaverUsecae) GetUserEmailByToken(accessToken string) (string, error) {
	const endpoint = "https://openapi.naver.com/v1/nid/me"

	// Create a new request to the specified endpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	// Set the Authorization header with the access token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var naverResp NaverResponse
	if err = json.Unmarshal(body, &naverResp); err != nil {
		return "", err
	}

	// Check if the message is not "success"
	if naverResp.Resultcode != "00" || naverResp.Message != "success" {
		return "", fmt.Errorf("failed to retrieve the email from Naver")
	}

	return naverResp.Response.Email, nil
}
