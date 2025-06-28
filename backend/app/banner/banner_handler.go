package banner

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"joosum-backend/pkg/util"
)

type BannerHandler struct {
	BannerUsecase BannerUsecase
}

// CreateBanner godoc
// @Summary 배너 생성
// @Description 새로운 배너를 생성합니다.
// @Tags 배너
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param request body banner.BannerCreateReq true "request"
// @Success 200 {object} Banner "배너 생성 성공"
// @Failure 400 {object} util.APIError "필수 파라미터가 누락된 경우 Bad Request를 반환합니다."
// @Failure 500 {object} util.APIError "배너 생성 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /banners [post]
func (h BannerHandler) CreateBanner(c *gin.Context) {

	var req BannerCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, util.MsgInvalidRequestBody)
		return
	}

	imageURL := req.ImageURL
	clickURL := req.ClickURL

	if imageURL == "" || clickURL == "" {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "필수 파라미터가 누락되었습니다")
		return
	}

	banner, err := h.BannerUsecase.CreateBanner(imageURL, clickURL)

	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, banner)
}

// GetBanners godoc
// @Summary 배너 목록 조회
// @Description 모든 배너 목록을 조회합니다.
// @Tags 배너
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} Banner "배너 목록 조회 성공"
// @Failure 500 {object} util.APIError "배너 목록 조회 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /banners [get]
func (h BannerHandler) GetBanners(c *gin.Context) {
	banners, err := h.BannerUsecase.GetBanners()

	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, banners)
}

// DeleteBanner godoc
// @Summary 배너 삭제
// @Description 배너를 삭제합니다.
// @Tags 배너
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param bannerId path string true "배너 ID"
// @Success 200 {object} int64 "배너 삭제 성공"
// @Failure 400 {object} util.APIError "필수 파라미터가 누락된 경우 Bad Request를 반환합니다."
// @Failure 500 {object} util.APIError "배너 삭제 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /banners/{bannerId} [delete]
func (h BannerHandler) DeleteBanner(c *gin.Context) {
	bannerId := c.Param("bannerId")

	if bannerId == "" {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "필수 파라미터가 누락되었습니다")
		return
	}

	count, err := h.BannerUsecase.DeleteBannerById(bannerId)

	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
