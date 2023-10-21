package auth

import (
	"fmt"
	"joosum-backend/app/link"
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase     AuthUsecase
	userUsecase     user.UserUsecase
	linkBookUsecase link.LinkBookUsecase
}

// SignUp godoc
// @Summary 회원 가입
// @Description 회원 가입을 위한 정보를 입력받고, 새로운 사용자를 생성하며 JWT 토큰 쌍을 반환합니다.
// @Tags 로그인
// @Accept  json
// @Produce  json
// @Param request body SignUpReq true "회원 가입 요청 본문"
// @Success 200 {object} util.TokenRes "회원 가입이 성공적으로 이루어지면 JWT 토큰 쌍을 반환합니다."
// @Failure 400 {object} util.APIError "요청 본문이 유효하지 않는 경우 Bad Request를 반환합니다."
// @Failure 409 {object} util.APIError "이미 존재하는 사용자의 경우 Conflict를 반환합니다."
// @Failure 500 {object} util.APIError "회원 가입 또는 JWT 토큰 생성 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /auth/signup [post]
func (h AuthHandler) SignUp(c *gin.Context) {
	var req SignUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIError{Error: "Invalid request body"})
		return
	}

	email, err := h.authUsecase.GetEmailFromJWT(req.Social, req.IdToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: fmt.Sprintf("failed to get email from the JWT token: %v", err.Error())})
		return
	}

	isExist, _ := h.userUsecase.GetUserByEmail(email)
	if isExist != nil {
		c.JSON(http.StatusConflict, util.APIError{Error: "user already exists"})
		return
	}

	inactiveUser, _ := h.userUsecase.GetInactiveUserByEmail(email)
	if inactiveUser != nil {
		c.JSON(http.StatusConflict, util.APIError{Error: "It hasn't been 30 days since you left"})
		return
	}

	temp_nickname := req.Nickname
	if temp_nickname == "" {
		temp_nickname = "user_" + util.RandomString(10)
	}

	var userInfo user.User

	userInfo.Email = email
	userInfo.Social = req.Social
	userInfo.Name = temp_nickname
	userInfo.Age = req.Age
	userInfo.Gender = req.Gender

	user, err := h.authUsecase.SignUp(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken(email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: err.Error()})
		return
	}

	// 회원가입 시 기본 링크북 폴더 생성
	h.linkBookUsecase.CreateDefaultLinkBook(user.UserId)

	c.JSON(http.StatusOK, util.TokenRes{AccessToken: accessToken, RefreshToken: refreshToken})
}

// Logout
// @Tags 유저
// @Summary 로그아웃
// @Description 디바이스 ID 삭제
// @Success 200 {object} db.UpdateResult
// @Security ApiKeyAuth
// @Router /auth/logout [POST]
func (h AuthHandler) Logout(c *gin.Context) {
	userId := util.GetUserId(c)

	result, err := h.authUsecase.Logout(userId)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, result)
}

// GetMe godoc
// @Summary 내 정보 조회
// @Description 현재 로그인된 사용자의 정보를 반환합니다.
// @Tags 로그인
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} user.User "현재 로그인된 사용자의 정보를 반환합니다."
// @Failure 401 {object} util.APIError "로그인이 되어있지 않은 경우 Unauthorized를 반환합니다."
// @Router /auth/me [get]
func (h AuthHandler) GetMe(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	userId := currentUser.(*user.User).UserId
	h.userUsecase.GetUserById(userId)

	c.JSON(http.StatusOK, currentUser)

}

// Protected
// @Tags 로그인
// @Summary 액세스토큰 테스트
// @Description 테스트하고자 하는 액세스토큰을 헤더에 넣고 요청을 보내면 success 를 반환합니다.
// @Security ApiKeyAuth
// @Router /protected [get]
func (h AuthHandler) Protected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

/*
curl -X 'POST' \
  'http://127.0.0.1:5001/auth/signup' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "access_token": "string",
  "age": 20,
  "email": "mono@test.com",
  "gender": "m",
  "nickname": "string",
  "social": "google"
}'
*/
