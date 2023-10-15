package auth

type NaverAccessTokenReq struct {
	IdToken string `json:"idToken"`
}

type NaverResponse struct {
	Resultcode string `json:"resultcode"`
	Message    string `json:"message"`
	Response   struct {
		Email string `json:"email"`
	} `json:"response"`
}
