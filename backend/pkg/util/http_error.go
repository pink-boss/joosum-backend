package util

type APIError struct {
	Error string `json:"error"`
}

type APIResponse struct {
	Message string `json:"message"`
}

type TokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
