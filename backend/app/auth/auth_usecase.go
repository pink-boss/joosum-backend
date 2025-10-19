package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"joosum-backend/app/setting"
	"joosum-backend/app/user"
	localConfig "joosum-backend/pkg/config"
	"joosum-backend/pkg/util"
	"math/big"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type AuthUsecase struct {
	salt         string
	userModel    user.UserModel
	settingModel setting.SettingModel
}

func (u *AuthUsecase) GenerateNewJWTToken(email string) (string, string, error) {
	accessToken, err := util.GenerateNewJWTAccessToken([]util.Role{util.User}, email)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := util.GenerateNewJWTRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *AuthUsecase) SignUp(userInfo user.User) (*user.User, error) {
	user, err := u.userModel.CreatUser(userInfo)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *AuthUsecase) Logout(userId string) (*mongo.UpdateResult, error) {
	result, err := u.settingModel.DeleteDivceId(userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *AuthUsecase) GetEmailFromJWT(social, idToken string) (string, error) {
	if social == "apple" {
		pubKey := ApplePublicKey{}
		publicSecret := PublicSecret{}
		client := resty.New()

		// 애플이 제공하는 public key 들을 가져옴
		pubKeyResult, err := client.R().SetResult(&pubKey).Get(ApplePublicKeyURL)
		JWKS := pubKeyResult.Result().(*ApplePublicKey)

		if err != nil {
			return "", fmt.Errorf("get apple public key err: %v", err)
		}

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(idToken, &claims, func(token *jwt.Token) (interface{}, error) {
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

		return email, nil
	} else if social == "google" {
		ctx := context.Background()

		// Google ID 토큰 검증 - 여러 플랫폼 지원 (iOS, Android, Web)
		audiences := []string{
			localConfig.GetEnvConfig("googleClientID"),        // iOS
			localConfig.GetEnvConfig("googleAndroidClientID"), // Android
			localConfig.GetEnvConfig("googleWebClientID"),     // Web
		}

		var lastErr error
		for _, audience := range audiences {
			oauth2Service, err := oauth2.NewService(ctx, option.WithAPIKey(localConfig.GetEnvConfig("googleApiKey")))
			if err != nil {
				lastErr = fmt.Errorf("unable to create OAuth2 service: %v", err)
				continue
			}

			tokenInfoCall := oauth2Service.Tokeninfo()
			tokenInfoCall.IdToken(idToken)

			tokenInfo, err := tokenInfoCall.Do()
			if err != nil {
				lastErr = fmt.Errorf("unable to verify IdToken with audience %s: %v", audience, err)
				continue
			}

			// Check if the token's audience matches the current audience
			if tokenInfo.Audience == audience {
				// Return the user's email address.
				if tokenInfo.Email != "" {
					return tokenInfo.Email, nil
				}
				return "", fmt.Errorf("unable to retrieve user's email from token")
			}
		}

		// 모든 audience로 검증 실패
		if lastErr != nil {
			return "", fmt.Errorf("failed to verify Google ID token for all audiences: %v", lastErr)
		}
		return "", fmt.Errorf("Google ID token audience does not match any configured client IDs")
	} else {
		return "", fmt.Errorf("invalid social name")
	}
}
