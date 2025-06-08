package link

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/errgo.v2/errors"
	"joosum-backend/pkg/util"
)

type LinkUsecase struct {
	linkModel     LinkModel
	linkBookModel LinkBookModel
}

func (u LinkUsecase) CreateLink(url string, title string, userId string, linkBookId string, thumbnailURL string, tags []string) (*Link, error) {

	// URL 이 http:// 혹은 https:// 로 시작하지 않으면 https:// 를 붙입니다.
	url = util.EnsureHTTPPrefix(url)

	// linkBookId 가 root 이거나 빈 스트링이라면 기본 폴더에 저장
	var linkBookName string
	if linkBookId == "root" || linkBookId == "" {
		defaultLinkBook, err := u.linkBookModel.GetDefaultLinkBook(userId)
		if err != nil {
			return nil, err
		}

		linkBookId = defaultLinkBook.LinkBookId
		linkBookName = defaultLinkBook.Title
	} else {
		linkBookData, err := u.linkBookModel.GetLinkBookById(linkBookId)
		if err != nil {
			return nil, err
		}

		linkBookName = linkBookData.Title
	}

	link, err := u.linkModel.CreateLink(url, title, userId, linkBookId, linkBookName, thumbnailURL, tags)
	if err != nil {
		return nil, err
	}

	// 링크북 최근 링크등록일 업데이트
	err = u.linkBookModel.UpdateLinkBookLastSavedAt(linkBookId)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (u LinkUsecase) Get9LinksByUserId(userId string) ([]*Link, error) {
	links, err := u.linkModel.Get9LinksByUserId(userId)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (u LinkUsecase) FindOneLinkByLinkId(linkId string) (*Link, error) {
	link, err := u.linkModel.GetOneLinkByLinkId(linkId)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (u LinkUsecase) FindAllLinksByUserId(userId string, sort string, order string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserId(userId, sort, order)
	if err != nil {
		return nil, err
	}

	// if links length 0 return []
	if len(links) == 0 {
		return []*Link{}, nil
	}

	return links, nil
}

func (u LinkUsecase) FindAllLinksByUserIdAndSearch(userId string, search string, sort string, order string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserIdAndSearch(userId, search, sort, order)
	if err != nil {
		return nil, err
	}

	// if links length 0 return []
	if len(links) == 0 {
		return []*Link{}, nil
	}

	return links, nil
}

func (u LinkUsecase) FindAllLinksByUserIdAndLinkBookId(userId string, linkBookId string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserIdAndLinkBookId(userId, linkBookId)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (u LinkUsecase) FindAllLinksByUserIdAndLinkBookIdAndSearch(userId string, linkBookId string, search string, sort string, order string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserIdAndLinkBookIdAndSearch(userId, linkBookId, search, sort, order)
	if err != nil {
		return nil, err
	}

	// if links length 0 return []
	if len(links) == 0 {
		return []*Link{}, nil
	}

	return links, nil
}

func (u LinkUsecase) DeleteOneByLinkId(linkId string) error {
	err := u.linkModel.DeleteOneByLinkId(linkId)
	if err != nil {
		return err
	}

	return nil
}

func (u LinkUsecase) DeleteAllLinks(userId string, linkIds []string) (int64, error) {
	if linkIds[0] == "all" {
		deletedCount, err := u.linkModel.DeleteAllLinksByUserId(userId)
		if err != nil {
			return 0, err
		}

		return deletedCount, nil

	} else if strings.HasPrefix(linkIds[0], "Link-") {
		deletedCount, err := u.linkModel.DeleteAllLinksByLinkIds(userId, linkIds)
		if err != nil {
			return 0, err
		}

		return deletedCount, nil
	} else {
		return 0, errors.New("Invalid query parameter")
	}
}

func (u LinkUsecase) DeleteAllLinksByLinkBookId(userId string, linkBookId string) error {
	_, err := u.linkModel.DeleteAllLinksByLinkBookId(userId, linkBookId)
	if err != nil {
		return err
	}

	return nil
}

func (u LinkUsecase) UpdateReadByLinkId(linkId string) error {
	err := u.linkModel.UpdateReadCountByLinkId(linkId)
	if err != nil {
		return err
	}

	return nil

}

func (u LinkUsecase) UpdateLinkBookIdByLinkId(linkId string, linkBookId string) error {
	// find LinkBook by linkBookId
	linkBookData, err := u.linkBookModel.GetLinkBookById(linkBookId)
	if err != nil {
		return err
	}

	linkBookName := linkBookData.Title

	err = u.linkModel.UpdateLinkBookIdAndTitleByLinkId(linkId, linkBookId, linkBookName)
	if err != nil {
		return err
	}

	return nil
}

func (u LinkUsecase) UpdateTitleAndUrlByLinkId(linkId string, url string, title string, thumbnailURL string, tags []string) (*Link, error) {
	// URL 이 http:// 혹은 https:// 로 시작하지 않으면 https:// 를 붙입니다.
	url = util.EnsureHTTPPrefix(url)
	link, err := u.linkModel.UpdateTitleAndUrlByLinkId(linkId, url, title, thumbnailURL, tags)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (LinkUsecase) GetThumnailURL(url string) (*LinkThumbnailRes, error) {
	// URL 이 http:// 혹은 https:// 로 시작하지 않으면 https:// 를 붙입니다.
	url = util.EnsureHTTPPrefix(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return &LinkThumbnailRes{
			URL:          url,
			ThumbnailURL: nil,
			Title:        &url,
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var ogTitle, ogImage, thumbnail *string

	doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if property, exists := item.Attr("property"); exists {
			if property == "og:title" {
				content, _ := item.Attr("content")
				ogTitle = &content
			}
			if property == "og:image" {
				content, _ := item.Attr("content")
				ogImage = &content

				if strings.HasPrefix(*ogImage, "//") {
					*ogImage = "https:" + *ogImage
				}
			}
		}
		if name, exists := item.Attr("name"); exists {
			if name == "thumbnail" {
				content, _ := item.Attr("content")
				thumbnail = &content
			}
		}
	})

	// meta tag 에서 없는 경우들 처리
	// og:title 이 없는 경우 title tag 에서 가져옴
	if ogTitle == nil {
		doc.Find("title").Each(func(index int, item *goquery.Selection) {
			ogTitle = new(string)
			*ogTitle = item.Text()
		})
	}

	// schema.org/SearchResultsPage를 사용하는 경우
	if ogImage == nil {
		// itemtype="http://schema.org/SearchResultsPage"를 가진 요소 찾기
		doc.Find(`[itemtype="http://schema.org/SearchResultsPage"]`).Each(func(index int, item *goquery.Selection) {
			// 필요한 데이터 추출 (예시: 썸네일 이미지)
			item.Find(`meta[itemprop="image"]`).Each(func(i int, img *goquery.Selection) {
				content, exists := img.Attr("content")
				if exists {
					ogImage = &content
				} else {
					src, exists := img.Attr("src")
					if exists {
						ogImage = &src
					}
				}
			})
		})
	}

	// google search에서 사용하는 경우
	if ogImage == nil && strings.Contains(url, "google.com/search") {
		// 무조건 https://www.gstatic.com/images/branding/googleg/1x/googleg_standard_color_128dp.png 사용
		ogImage = new(string)
		*ogImage = "https://www.gstatic.com/images/branding/googleg/1x/googleg_standard_color_128dp.png"
	}

	// Attempt to find the first valid image URL
	if ogImage == nil {
		doc.Find("img").Each(func(index int, item *goquery.Selection) {
			src, exists := item.Attr("src")
			if exists && !strings.HasPrefix(src, "data:image") { // Ignore base64 images
				ogImage = &src
				return
			}
		})
	}

	// twitter:image를 사용하는 경우
	if ogImage == nil {
		doc.Find(`meta[name="twitter:image"]`).Each(func(index int, item *goquery.Selection) {
			content, exists := item.Attr("content")
			if exists {
				ogImage = &content
			}
		})
	}

	// twitter:title를 사용하는 경우
	if ogTitle == nil {
		doc.Find(`meta[name="twitter:title"]`).Each(func(index int, item *goquery.Selection) {
			content, exists := item.Attr("content")
			if exists {
				ogTitle = &content
			}
		})
	}

	// Thumbnail URL이 ogImage보다 우선순위가 낮기 때문에 마지막에 설정
	if ogImage == nil && thumbnail != nil {
		ogImage = thumbnail
	}

	// 만약 http 나 https로 시작하지 않는 경우, null 반환
	if ogImage != nil && !strings.HasPrefix(*ogImage, "http") {
		ogImage = nil
	}

	return &LinkThumbnailRes{
		URL:          url,
		ThumbnailURL: ogImage,
		Title:        ogTitle,
	}, nil

}
