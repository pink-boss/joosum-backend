package util

import (
	"github.com/go-resty/resty/v2"
)

var JWKS *ApplePublicKey

const ApplePublicKeyURL = "https://appleid.apple.com/auth/keys"

func LoadApplePublicKeys() {
	pubKey := ApplePublicKey{}
	client := resty.New()

	// 애플이 제공하는 public key 들을 가져옴
	pubKeyResult, err := client.R().SetResult(&pubKey).Get(ApplePublicKeyURL)
	JWKS = pubKeyResult.Result().(*ApplePublicKey)

	if err != nil {
		panic("failed to get public keys from the apple")
	}
}

type appleKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"Use"`
	Als string `json:"Als"`
	N   string `json:"N"`
	E   string `json:"E"`
}

type PublicSecret struct {
	N []byte
	E []byte
}

type ApplePublicKey struct {
	Keys []appleKey `json:"Keys"`
}
