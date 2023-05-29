package auth

import (
	"net/http"

	"joosum-backend/app/user"
	"joosum-backend/pkg/util"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase *AuthUsecase
	userUsecase *user.UserUsecase
}

// SignUp godoc
// @Summary 회원 가입
// @Description 회원 가입을 위한 정보를 입력받고, 새로운 사용자를 생성하며 JWT 토큰 쌍을 반환합니다.
// @Tags 로그인
// @Accept  json
// @Produce  json
// @Param request body SignUpRequest true "회원 가입 요청 본문"
// @Success 200 {object} util.TokenResponse "회원 가입이 성공적으로 이루어지면 JWT 토큰 쌍을 반환합니다."
// @Failure 400 {object} util.APIError "요청 본문이 유효하지 않는 경우 Bad Request를 반환합니다."
// @Failure 409 {object} util.APIError "이미 존재하는 사용자의 경우 Conflict를 반환합니다."
// @Failure 500 {object} util.APIError "회원 가입 또는 JWT 토큰 생성 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIError{Error: "Invalid request body"})
		return
	}

	isExist, err := h.userUsecase.GetUserByEmail(req.Email); 
	if isExist != nil {
		c.JSON(http.StatusConflict, util.APIError{Error: "user already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: err.Error()})
		return
	}

	var userInfo user.User
	userInfo.Email = req.Email
	userInfo.Social = req.Social
	userInfo.Name = req.Nickname
	userInfo.Age = req.Age
	userInfo.Gender = req.Gender

	_, err = h.authUsecase.SignUp(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken([]string{"USER", "ADMIN"}, req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.TokenResponse{AccessToken: accessToken , RefreshToken: refreshToken})
}
