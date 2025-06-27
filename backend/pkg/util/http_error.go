package util

import "github.com/gin-gonic/gin"

// APIError 는 API에서 발생하는 오류 응답 형식을 정의합니다.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIResponse 는 일반적인 메시지 응답 형식을 정의합니다.
type APIResponse struct {
	Message string `json:"message"`
}

// TokenRes 는 액세스 토큰과 리프레시 토큰 쌍을 반환할 때 사용됩니다.
type TokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// 사전 정의된 오류 코드
const (
	CodeInvalidRequestBody   = "INVALID_REQUEST_BODY"
	CodeMissingAuthorization = "MISSING_AUTHORIZATION"
	CodeInternalServerError  = "INTERNAL_SERVER_ERROR"
)

// 사전 정의된 오류 메시지(한글)
const (
	MsgInvalidRequestBody   = "잘못된 요청 본문입니다."
	MsgMissingAuthorization = "Authorization 헤더가 없습니다."
	MsgInternalServerError  = "서버 오류가 발생했습니다."
)

// SendError 는 오류 응답을 JSON 형태로 클라이언트에 반환합니다.
func SendError(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIError{Code: code, Message: message})
}
