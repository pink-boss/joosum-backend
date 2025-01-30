package auth

import (
	"fmt"
	"joosum-backend/app/user"
	localConfig "joosum-backend/pkg/config"
	"joosum-backend/pkg/util"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gin-gonic/gin"
)

// https://gamz.tistory.com/11
type GoogleHandler struct {
	// DI
	// please write private
	authUsecase  AuthUsecase
	googleUsecae GoogleUsecae
	userUsecase  user.UserUsecase
}

// @Summary Google 액세스 토큰 검증
// @Description Google 액세스 토큰의 유효성을 검사하고 새 JWT 토큰 쌍을 생성합니다.
// @Tags 로그인
// @Accept  json
// @Produce  json
// @Param request body AccessTokenReq true "액세스 토큰 요청 본문"
// @Success 200 {object} util.TokenRes "액세스 토큰을 성공적으로 검증하고 JWT 토큰을 반환합니다."
// @Failure 400 {object} util.APIError "요청 본문이 유효하지 않거나 액세스토큰이 누락된 경우 Bad Request를 반환합니다."
// @Failure 401 {object} util.APIError "Google 액세스 토큰이 유효하지 않거나 사용자가 존재하지 않는 경우 Unauthorized를 반환합니다."
// @Failure 500 {object} util.APIError "액세스 토큰을 검증하거나 JWT 토큰을 생성하는 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /auth/google [post]
func (h *GoogleHandler) VerifyGoogleAccessToken(c *gin.Context) {
	var req AccessTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIError{Error: "Invalid request body"})
		return
	}

	accessToken := req.IdToken

	if accessToken == "" {
		c.JSON(http.StatusBadRequest, util.APIError{Error: "idToken is required"})
		return
	}

	if valid, err := h.googleUsecae.ValidateIdToken(accessToken); err != nil || !valid {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: "Invalid id token"})
		return
	}

	email, err := h.googleUsecae.GetUserEmail(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.APIError{Error: err.Error()})
		return
	}

	user, err := h.userUsecase.GetUserByEmail(email)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("failed to get user by email: %v", err.Error())})
		return
	}

	if user == nil {
		c.JSON(http.StatusOK, util.TokenRes{AccessToken: "", RefreshToken: ""})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken(email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.TokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken})

}

// 웹/안드로이드용 OAuth2 설정 (google_handler.go 안에서 바로 선언)
var googleWebOAuthConfig = &oauth2.Config{
	ClientID:     localConfig.GetEnvConfig("googleWebClientID"),
	ClientSecret: localConfig.GetEnvConfig("googleWebSecret"),
	RedirectURL:  localConfig.GetEnvConfig("googleWebRedirect"), // ex) https://yourdomain.com/auth/google/callback
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

// 2.1) 구글 로그인 페이지로 리다이렉트
func (h *GoogleHandler) AuthGoogleWebLogin(c *gin.Context) {
	state := util.GenerateRandomState(16) // CSRF 방지용 랜덤 state
	// state를 쿠키나 세션에 저장 (쿠키 예시)
	c.SetCookie("google_oauth_state", state, 3600, "/", "", false, true)

	// Google 로그인 URL 생성
	url := googleWebOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// 2.2) 구글 인증 후 서버로 돌아오는 콜백
func (h *GoogleHandler) AuthGoogleWebCallback(c *gin.Context) {
	// 1) state 검증
	state := c.Query("state")
	storedState, err := c.Cookie("google_oauth_state")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state cookie not found"})
		return
	}
	if state != storedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	// 2) code 가져오기
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
		return
	}

	// 3) code로 토큰 교환
	token, err := googleWebOAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to exchange token: " + err.Error()})
		return
	}

	// 4) 토큰으로 구글 사용자 정보 가져오기
	userInfo, err := h.googleUsecae.GetUserInfoFromToken(c, token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info: " + err.Error()})
		return
	}

	email := userInfo.Email
	if email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no email in user info"})
		return
	}

	// 5) DB에서 해당 이메일 유저 조회
	foundUser, err := h.userUsecase.GetUserByEmail(email)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user: " + err.Error()})
		return
	}

	// 5-1) 유저가 없으면 회원가입 로직 (예시로 빈 토큰 리턴)
	if foundUser == nil {
		c.JSON(http.StatusOK, gin.H{"message": "user not found, please sign up first"})
		return
	}

	// 6) JWT 토큰 발급
	accessTokenStr, refreshTokenStr, err := h.authUsecase.GenerateNewJWTToken(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT: " + err.Error()})
		return
	}

	// 7) 응답 (AccessToken/RefreshToken)
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessTokenStr,
		"refreshToken": refreshTokenStr,
	})
}
