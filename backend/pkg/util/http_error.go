package util

type APIError struct {
    Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}