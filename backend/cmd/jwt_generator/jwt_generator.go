package main

import (
	"fmt"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/util"
)

func main() {
	config.EnvConfig()

	token, err := util.GenerateNewJWTAccessToken([]string{"USER", "ADMIN"}, "admin@gmail.com")
	if err != nil {
		println(err)
	}
	fmt.Println(token)
}
