package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/util"
	"math/big"
	"time"
)

const ApplePublicKey = "https://appleid.apple.com/auth/keys"
const AppleBaseURL = "https://appleid.apple.com"

func issueTokenFromApple(reqAuth authRequest) (*clientResponse, error) {
	pubKey := applePublicKey{}
	publicSecret := publicSecret{}
	client := resty.New()

	// 애플이 제공하는 public key 들을 가져옴
	pubKeyResult, err := client.R().SetResult(&pubKey).Get(ApplePublicKey)
	JWKS := pubKeyResult.Result().(*applePublicKey)

	if err != nil {
		return nil, fmt.Errorf("get apple public key err: %v", err)
	}

	idTk, err := jwt.Parse(reqAuth.IdToken, func(token *jwt.Token) (interface{}, error) {
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
	if err != nil {
		return nil, fmt.Errorf("unable to verify access token: %v", err)
	}

	// verify 가 완료된 idTk의 payload 값을 받아와서 맞는 정보인지 검증
	claims, ok := idTk.Claims.(jwt.MapClaims)
	if !ok && !idTk.Valid {
		return nil, fmt.Errorf("token is not valid: %v", err)
	}

	// todo 사용자 정보 및 로그인 히스토리 DB 에 저장
	fmt.Println(claims)

	// 애플한테 access, refresh token 받기

	clientID := config.GetEnvConfig("apple.clientID")
	teamID := config.GetEnvConfig("apple.teamID")
	keyID := config.GetEnvConfig("apple.keyID")

	appleClaims := jwt.MapClaims{
		"iss": teamID,
		"aud": AppleBaseURL,
		"exp": time.Now().UTC().Add(24 * time.Hour * 100).Unix(),
		"iat": time.Now().UTC().Unix(),
		"sub": clientID,
	}

	appleToken := jwt.NewWithClaims(jwt.SigningMethodES256, appleClaims)
	appleToken.Header["kid"] = keyID

	privateKey, err := util.LoadPrivateKey("Apple_AuthKey.p8")
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %v", err)
	}

	signedToken, err := appleToken.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("fail to signing with private key: %v", err)
	}

	formData := map[string]string{
		"client_id":     clientID,
		"client_secret": signedToken,
		"code":          reqAuth.Code,
		"grant_type":    "authorization_code",
	}

	resToken := clientResponse{}
	uri := AppleBaseURL + "/auth/token"
	result, err := client.R().SetFormData(formData).SetResult(&resToken).Post(uri)

	if result.IsError() {
		return nil, fmt.Errorf("fail to get the token from apple: %v", result.RawResponse)
	}
	if err != nil {
		return nil, fmt.Errorf("response get failure.: %v", err)
	}

	return result.Result().(*clientResponse), nil
}
