package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"joosum-backend/pkg/util"
	"math/big"
	"net/http"
	"strings"
)

func AppleAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		publicSecret := util.PublicSecret{}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		idTk, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// kid 값 저장 | public key 대조에 필요하기 때문에
			kid := token.Header["kid"].(string)

			for _, v := range util.JWKS.Keys {
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
		if err != nil || !idTk.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("invalid or expired token: %v", err.Error()),
			})
			c.Abort()
			return
		}

		// Token is valid, proceed with the next handler
		c.Next()
	}
}
