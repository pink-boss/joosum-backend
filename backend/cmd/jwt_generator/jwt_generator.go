package main

import (
	"fmt"
	"joosum-backend/pkg/util"
)

func main() {
	token, err := util.GenerateNewJWTAccessToken([]string{"USER", "ADMIN"}, "1", "fh6Bs8C")
	if err != nil {
		println(err)
	}
	fmt.Println(token)
}
