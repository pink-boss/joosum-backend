package auth

type authRequest struct {
	Code    string
	IdToken string `json:"id_token"`
}

//	{
//	 "access_token": "adg61...67Or9",
//	 "token_type": "Bearer",
//	 "expires_in": 3600,
//	 "refresh_token": "rca7...lABoQ",
//	 "id_token": "eyJra...96sZg"
//	}
type clientResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token" example:"adg61...67Or9"`
	RefreshToken string `json:"refresh_token" example:"rca7...lABoQ"`
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
