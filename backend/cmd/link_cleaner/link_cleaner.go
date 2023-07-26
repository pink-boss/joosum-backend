package main

import (
	"fmt"
	"joosum-backend/app/link"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/util"
)

// 링크북이 없는 고아 링크를 삭제
func main() {
	config.EnvConfig()
	util.StartMongoDB()

	linkModel := link.LinkModel{}
	linkBookModel := link.LinkBookModel{}

	links, _ := linkModel.GetAllLink()
	deletedCount := 0

	for _, link := range links {
		haveParent := linkBookModel.HaveLinkBook(link.LinkBookId)

		// 부모 링크북이 없으면 링크 삭제
		if haveParent == false {
			linkModel.DeleteOneByLinkId(link.LinkId)
			deletedCount += 1
		}
	}

	fmt.Printf("링크북이 없는 %d 개의 링크가 삭제되었습니다.", deletedCount)
}
