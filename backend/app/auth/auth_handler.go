package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase AuthUsecase
	userUsecase user.UserUsecase
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
func (h AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIError{Error: "Invalid request body"})
		return
	}

	isExist, _ := h.userUsecase.GetUserByEmail(req.Email)
	if isExist != nil {
		c.JSON(http.StatusConflict, util.APIError{Error: "user already exists"})
		return
	}

	var userInfo user.User

	if req.Social == "apple" {
		pubKey := ApplePublicKey{}
		publicSecret := PublicSecret{}
		client := resty.New()

		// 애플이 제공하는 public key 들을 가져옴
		pubKeyResult, err := client.R().SetResult(&pubKey).Get(ApplePublicKeyURL)
		JWKS := pubKeyResult.Result().(*ApplePublicKey)

		if err != nil {
			//return nil, fmt.Errorf("get apple public key err: %v", err)
			return
		}

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(req.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
			// kid 값 저장 | public key 대조에 필요하기 때문에
			kid := token.Header["kid"].(string)

			for _, v := range JWKS.Keys {
				// 받아온 public key 중 id_token 과 kid 일치하는지 확인 후 n, e 값 저장
				if kid == v.Kid {
					n, _ := base64.RawURLEncoding.DecodeString(v.N)
					publicSecret.N = n
					e, _ := base64.StdEncoding.DecodeString(v.E)
					publicSecret.E = e
					break
				}
			}
			publicKeyE := binary.LittleEndian.Uint32(append(publicSecret.E, 0))

			// 이 rsaKey 로 id_token verify
			rsaKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(publicSecret.N),
				E: int(publicKeyE),
			}
			return rsaKey, nil
		})

		email := claims["email"].(string)
		userInfo.Email = email
	} else {
		userInfo.Email = req.Email
	}

	userInfo.Social = req.Social
	userInfo.Name = req.Nickname
	userInfo.Age = req.Age
	userInfo.Gender = req.Gender

	_, err := h.authUsecase.SignUp(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken([]string{"USER", "ADMIN"}, req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
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
