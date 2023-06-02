package auth

type SignUpRequest struct {
	AccessToken string `json:"accessToken"`
	Social      string `json:"social"  example:"google"`
	Nickname    string `json:"nickname"`
	Gender      string `json:"gender" example:"m"`
	Age         uint8  `json:"age" example:"20"`
}

type SignUpInfo struct {
	Email      string `json:"email"`
	Social     string `json:"social"  example:"google"`
	Nickname   string `json:"nickname"`
	Gender     string `json:"Gender" example:"m"`
	Age        uint8  `json:"age" example:"20"`
	SignUpDate string `json:"sign_up_date"`
}
