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

// @Summary (apple)Google 액세스 토큰 검증
// @Description (apple)Google 액세스 토큰의 유효성을 검사하고 새 JWT 토큰 쌍을 생성합니다.
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
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, util.MsgInvalidRequestBody)
		return
	}

	accessToken := req.IdToken

	if accessToken == "" {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "idToken 값이 필요합니다")
		return
	}

	if valid, err := h.googleUsecae.ValidateIdToken(accessToken); err != nil || !valid {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, "유효하지 않은 idToken")
		return
	}

	email, err := h.googleUsecae.GetUserEmail(accessToken)
	if err != nil {
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, err.Error())
		return
	}

	user, err := h.userUsecase.GetUserByEmail(email)
	if err != nil && err != mongo.ErrNoDocuments {
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, fmt.Sprintf("failed to get user by email: %v", err.Error()))
		return
	}

	if user == nil {
		c.JSON(http.StatusOK, util.TokenRes{AccessToken: "", RefreshToken: ""})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken(email)

	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, util.TokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken})

}

// @Summary (android)Google 액세스 토큰 검증
// @Description (android)Google 액세스 토큰의 유효성을 검사하고 새 JWT 토큰 쌍을 생성합니다.
// @Tags 로그인
// @Accept  json
// @Produce  json
// @Param request body AccessTokenReq true "액세스 토큰 요청 본문"
// @Success 200 {object} util.TokenRes "액세스 토큰을 성공적으로 검증하고 JWT 토큰을 반환합니다."
// @Failure 400 {object} util.APIError "요청 본문이 유효하지 않거나 액세스토큰이 누락된 경우 Bad Request를 반환합니다."
// @Failure 401 {object} util.APIError "Google 액세스 토큰이 유효하지 않거나 사용자가 존재하지 않는 경우 Unauthorized를 반환합니다."
// @Failure 500 {object} util.APIError "액세스 토큰을 검증하거나 JWT 토큰을 생성하는 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /auth/google [post]
func (h *GoogleHandler) VerifyGoogleAccessTokenInAndroid(c *gin.Context) {
	var req AccessTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, util.MsgInvalidRequestBody)
		return
	}

	accessToken := req.IdToken

	if accessToken == "" {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "idToken 값이 필요합니다")
		return
	}

	if valid, err := h.googleUsecae.ValidateIdTokenForAndroid(accessToken); err != nil || !valid {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, "유효하지 않은 idToken")
		return
	}

	email, err := h.googleUsecae.GetUserEmail(accessToken)
	if err != nil {
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, err.Error())
		return
	}

	user, err := h.userUsecase.GetUserByEmail(email)
	if err != nil && err != mongo.ErrNoDocuments {
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, fmt.Sprintf("failed to get user by email: %v", err.Error()))
		return
	}

	if user == nil {
		c.JSON(http.StatusOK, util.TokenRes{AccessToken: "", RefreshToken: ""})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken(email)

	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, util.TokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken})

}

func (h *GoogleHandler) VerifyGoogleAccessTokenInWeb(c *gin.Context) {
	var req AccessTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, util.MsgInvalidRequestBody)
		return
	}

	accessToken := req.IdToken

	if accessToken == "" {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "idToken 값이 필요합니다")
		return
	}

	valid, err := h.googleUsecae.ValidateIdTokenForWeb(accessToken)
	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, fmt.Sprintf("Invalid id token: %v", err))
		return
	}

	if !valid {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, "유효하지 않은 idToken")
		return
	}

	email, err := h.googleUsecae.GetUserEmail(accessToken)
	if err != nil {
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, err.Error())
		return
	}

	user, err := h.userUsecase.GetUserByEmail(email)
	if err != nil && err != mongo.ErrNoDocuments {
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, fmt.Sprintf("failed to get user by email: %v", err.Error()))
		return
	}

	if user == nil {
		c.JSON(http.StatusOK, util.TokenRes{AccessToken: "", RefreshToken: ""})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken(email)

	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, util.TokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
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
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "state 쿠키를 찾을 수 없습니다")
		return
	}
	if state != storedState {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "잘못된 OAuth 상태입니다")
		return
	}

	// 2) code 가져오기
	code := c.Query("code")
	if code == "" {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "코드를 찾을 수 없습니다")
		return
	}

	// 3) code로 토큰 교환
	token, err := googleWebOAuthConfig.Exchange(c, code)
	if err != nil {
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, "토큰 교환 실패: "+err.Error())
		return
	}

	// 4) 토큰으로 구글 사용자 정보 가져오기
	userInfo, err := h.googleUsecae.GetUserInfoFromToken(c, token.AccessToken)
	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, "유저 정보를 가져오지 못했습니다: "+err.Error())
		return
	}

	email := userInfo.Email
	if email == "" {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, "유저 정보에 이메일이 없습니다")
		return
	}

	// 5) DB에서 해당 이메일 유저 조회
	foundUser, err := h.userUsecase.GetUserByEmail(email)
	if err != nil && err != mongo.ErrNoDocuments {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, "유저 조회 실패: "+err.Error())
		return
	}

	// 5-1) 유저가 없으면 회원가입 로직 (예시로 빈 토큰 리턴)
	if foundUser == nil {
		c.JSON(http.StatusOK, util.APIResponse{Message: "사용자를 찾을 수 없습니다. 먼저 회원가입을 진행해 주세요"})
		return
	}

	// 6) JWT 토큰 발급
	accessTokenStr, refreshTokenStr, err := h.authUsecase.GenerateNewJWTToken(email)
	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, "JWT 생성 실패: "+err.Error())
		return
	}

	// 7) 응답 (AccessToken/RefreshToken)
	c.JSON(http.StatusOK, util.TokenRes{AccessToken: accessTokenStr, RefreshToken: refreshTokenStr})
}
