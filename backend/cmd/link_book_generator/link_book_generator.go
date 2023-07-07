package main

import (
	"fmt"
	"joosum-backend/app/link"
	userPkg "joosum-backend/app/user"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/util"
)

func main() {
	config.EnvConfig()
	util.StartMongoDB()

	userModel := userPkg.UserModel{}
	linkBookModel := link.LinkBookModel{}
	linkBookUsecase := link.LinkBookUsecase{}

	users, _ := userModel.FindUsers()
	CreatedCount := 0

	// 모든 회원 검색
	for _, user := range users {
		linkBook, _ := linkBookModel.GetDefaultLinkBook(user.UserId)

		// 기본 링크북 폴더가 없으면 생성
		if linkBook == nil {
			_, err := linkBookUsecase.CreateDefaultLinkBook(user.UserId)
			if err == nil {
				CreatedCount += 1
			}
		}
	}

	fmt.Printf("%d 개의 기본 링크북 폴더가 생성되었습니다.", CreatedCount)
}
