package auth

type authRequest struct {
	State   string
	Code    string
	IdToken string `json:"id_token"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type appleKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"Use"`
	Als string `json:"Als"`
	N   string `json:"N"`
	E   string `json:"E"`
}

type publicSecret struct {
	N []byte
	E []byte
}

type applePublicKey struct {
	Keys []appleKey `json:"Keys"`
}
