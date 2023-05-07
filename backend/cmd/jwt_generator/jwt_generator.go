package main

import (
	"fmt"
	"joosum-backend/pkg/util"
)

func main() {
	token, err := util.GenerateNewJWTAccessToken([]string{"USER", "ADMIN"}, "fh6Bs8C")
	if err != nil {
		println(err)
	}
	fmt.Println(token)
}
