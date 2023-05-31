package util

type APIError struct {
    Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}