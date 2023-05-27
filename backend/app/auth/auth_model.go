package auth

type SignUpRequest struct {
	Email  string `json:"email"`
	Social string `json:"social"`
}
