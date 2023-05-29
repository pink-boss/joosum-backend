package auth

type authRequest struct {
	IdToken string `json:"idToken" example:"eyJra...LFmZQ"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken" example:"eyJhb...S_HK4"`
	RefreshToken string `json:"refreshToken" example:"46c67...2f891"`
}

type appleKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"Use"`
	Als string `json:"Als"`
	N   string `json:"N"`
	E   string `json:"E"`
}

type PublicSecret struct {
	N []byte
	E []byte
}

type ApplePublicKey struct {
	Keys []appleKey `json:"Keys"`
}
