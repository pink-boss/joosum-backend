package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"math/big"
)

const ApplePublicKey = "https://appleid.apple.com/auth/keys"

type appleKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"Use"`
	Als string `json:"Als"`
	N   string `json:"N"`
	E   string `json:"E"`
}

type publicSecret struct {
	N []byte
	E []byte
}

type applePublicKey struct {
	Keys []appleKey `json:"Keys"`
}

func verifyToken(reqAuth authRequest) (jwt.MapClaims, error) {
	pubKey := applePublicKey{}
	publicSecret := publicSecret{}
	client := resty.New()

	// 애플이 제공하는 public key 들을 가져옴
	pubKeyResult, err := client.R().SetResult(&pubKey).Get(ApplePublicKey)
	result := pubKeyResult.Result().(*applePublicKey)

	if err != nil {
		return nil, fmt.Errorf("get apple public key err")
	}

	idTk, err := jwt.Parse(reqAuth.IdToken, func(token *jwt.Token) (interface{}, error) {
		// kid 값 저장 | public key 대조에 필요하기 때문에
		kid := token.Header["kid"].(string)

		for _, v := range result.Keys {
			// 받아온 public key중 id_token과 kid 일치하는지 확인 후 n, e 값 저장
			if kid == v.Kid {
				n, _ := base64.RawURLEncoding.DecodeString(v.N)
				publicSecret.N = n
				e, _ := base64.StdEncoding.DecodeString(v.E)
				publicSecret.E = e
				break
			}
		}
		publicKeyE := binary.LittleEndian.Uint32(append(publicSecret.E, 0))

		// 이 rsaKey로 id_token verify
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

	return claims, nil
}
