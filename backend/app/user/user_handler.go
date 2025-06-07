package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase UserUsecase
}

// GetMainPage
// @Tags Main
// @Summary
// @Router / [get]
func GetMainPage(c *gin.Context) {
	log.Println("Main page....")
	c.String(http.StatusOK, "Main page for secure API!!")
}

// DeleteUser
// @Tags 유저
// @Summary 유저를 삭제하고 링크, 링크북을 삭제한다. 비활성유저로 변경된다.
// @Description
// @Accept json
// @Produce json
// @Success 200 {object} util.APIResponse "유저 삭제 성공"
// @Failure 401 {object} util.APIError "Authorization 헤더가 없을 때 반환합니다."
// @Failure 404 {object} util.APIError "유저를 찾을 수 없음"
// @Failure 500 {object} util.APIError "서버에서 유저 삭제 실패"
// @Router /auth/me [delete]
func (h UserHandler) DeleteUser(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}
	email := currentUser.(*User).Email

	err := h.userUsecase.UpdateUserToInactiveUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버에서 유저 삭제 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "유저 삭제 성공"})
}

// GetWithdrawUsers
// @Tags 유저
// @Summary 탈퇴한 유저 목록을 조회합니다.
// @Description
// @Accept json
// @Produce json
// @Success 200 {array} InactiveUser
// @Failure 500 {object} util.APIError "서버에서 유저 조회 실패"
// @Router /withdraw-users [get]
func (h UserHandler) GetWithdrawUsers(c *gin.Context) {
	users, err := h.userUsecase.GetWithdrawUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버에서 유저 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, users)
}
