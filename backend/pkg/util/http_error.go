package util

import "github.com/gin-gonic/gin"

// APIError 는 API에서 발생하는 오류 응답 형식을 정의합니다.
type APIError struct {
	Code    int    `json:"code" example:"1000"`
	Message string `json:"message" example:"잘못된 요청 본문입니다."`
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
	CodeInvalidRequestBody   = 1000
	CodeMissingAuthorization = 1001
	CodeInternalServerError  = 1002
	CodeMissingParameter     = 1003

	CodeInvalidIDToken   = 2000
	CodeUserExists       = 2001
	CodeUserRecentlyLeft = 2002
	CodeDuplicateTitle   = 3000
)

// 사전 정의된 오류 메시지(한글)
var codeMessages = map[int]string{
	CodeInvalidRequestBody:   "잘못된 요청 본문입니다.",
	CodeMissingAuthorization: "Authorization 헤더가 없습니다.",
	CodeInternalServerError:  "서버 오류가 발생했습니다.",
	CodeMissingParameter:     "필수 파라미터가 누락되었습니다.",
	CodeInvalidIDToken:       "유효하지 않은 ID 토큰입니다.",
	CodeUserExists:           "이미 존재하는 사용자입니다.",
	CodeUserRecentlyLeft:     "탈퇴 후 30일이 지나지 않았습니다.",
	CodeDuplicateTitle:       "같은 이름의 폴더가 존재합니다.",
}

// SendError 는 오류 응답을 JSON 형태로 클라이언트에 반환합니다.
func SendError(c *gin.Context, status int, code int) {
	msg, ok := codeMessages[code]
	if !ok {
		msg = "알 수 없는 오류가 발생했습니다."
	}
	c.JSON(status, APIError{Code: code, Message: msg})
}

// ErrorMessage 는 코드에 해당하는 메시지를 반환합니다.
func ErrorMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return ""
}
