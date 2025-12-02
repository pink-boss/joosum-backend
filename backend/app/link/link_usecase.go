package link

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/errgo.v2/errors"
	localConfig "joosum-backend/pkg/config"
	"joosum-backend/pkg/util"
)

// YouTubeVideoInfo YouTube API에서 가져온 동영상 정보
type YouTubeVideoInfo struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	CategoryID  string   `json:"categoryId"`
}

// YouTubeAPIResponse YouTube Data API v3 응답 구조
type YouTubeAPIResponse struct {
	Items []struct {
		Snippet struct {
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Tags        []string `json:"tags"`
			CategoryID  string   `json:"categoryId"`
		} `json:"snippet"`
	} `json:"items"`
}

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

// extractYouTubeVideoID YouTube URL에서 video ID를 추출합니다
func extractYouTubeVideoID(urlStr string) string {
	// youtube.com/watch?v=VIDEO_ID 형식
	if strings.Contains(urlStr, "youtube.com/watch") {
		parsedURL, err := url.Parse(urlStr)
		if err == nil {
			return parsedURL.Query().Get("v")
		}
	}

	// youtu.be/VIDEO_ID 형식
	if strings.Contains(urlStr, "youtu.be/") {
		parts := strings.Split(urlStr, "youtu.be/")
		if len(parts) > 1 {
			videoID := strings.Split(parts[1], "?")[0]
			return strings.TrimSpace(videoID)
		}
	}

	return ""
}

// getYouTubeVideoInfo YouTube Data API를 사용하여 동영상 정보를 가져옵니다
func getYouTubeVideoInfo(videoID string) (*YouTubeVideoInfo, error) {
	apiKey := localConfig.GetEnvConfig("youtubeApiKey")
	if apiKey == "" {
		log.Printf("[YouTube API] API 키가 설정되지 않았습니다. 크롤링 방식으로 진행합니다.")
		return nil, nil // API 키가 없으면 nil 반환 (에러 아님)
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s", videoID, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("[YouTube API] API 호출 실패 (videoID=%s): %v", videoID, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[YouTube API] API 응답 코드 비정상 (videoID=%s, status=%d)", videoID, resp.StatusCode)
		return nil, fmt.Errorf("YouTube API returned status %d", resp.StatusCode)
	}

	var apiResponse YouTubeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("[YouTube API] 응답 파싱 실패 (videoID=%s): %v", videoID, err)
		return nil, err
	}

	if len(apiResponse.Items) == 0 {
		log.Printf("[YouTube API] 동영상을 찾을 수 없습니다 (videoID=%s)", videoID)
		return nil, fmt.Errorf("video not found")
	}

	snippet := apiResponse.Items[0].Snippet
	videoInfo := &YouTubeVideoInfo{
		Title:       snippet.Title,
		Description: snippet.Description,
		Tags:        snippet.Tags,
		CategoryID:  snippet.CategoryID,
	}

	log.Printf("[YouTube API] 동영상 정보 가져오기 성공 (videoID=%s, title=%s, tags=%d개)",
		videoID, videoInfo.Title, len(videoInfo.Tags))

	return videoInfo, nil
}

// GetAIRecommendedTags AI를 사용하여 URL의 본문 내용을 분석하고 추천 태그를 생성합니다
func (LinkUsecase) GetAIRecommendedTags(url string) (*AITagRecommendationRes, error) {
	// URL 이 http:// 혹은 https:// 로 시작하지 않으면 https:// 를 붙입니다.
	url = util.EnsureHTTPPrefix(url)

	// YouTube URL 감지
	isYouTube := strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be")

	// User-Agent와 헤더를 설정한 HTTP 클라이언트 생성
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[AI 태그 추천] HTTP 요청 생성 실패 (url=%s): %v", url, err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 브라우저처럼 보이도록 헤더 설정 (봇 차단 방지)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7")
	// Accept-Encoding은 설정하지 않음 - Go의 http.Client가 자동으로 gzip 처리

	// URL에서 본문 내용 크롤링
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[AI 태그 추천] URL 크롤링 실패 (url=%s): %v", url, err)
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[AI 태그 추천] URL 응답 코드 비정상 (url=%s, status=%s)", url, resp.Status)
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("[AI 태그 추천] HTML 파싱 실패 (url=%s): %v", url, err)
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// 본문 텍스트 추출
	var contentBuilder strings.Builder

	// 1. 제목 추출 (우선순위: og:title > twitter:title > title 태그)
	title := ""
	doc.Find("meta[property='og:title']").Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists && title == "" {
			title = content
		}
	})
	if title == "" {
		doc.Find("meta[name='twitter:title']").Each(func(i int, s *goquery.Selection) {
			if content, exists := s.Attr("content"); exists {
				title = content
			}
		})
	}
	if title == "" {
		title = doc.Find("title").First().Text()
	}
	title = strings.TrimSpace(title)
	if title != "" {
		contentBuilder.WriteString("제목: ")
		contentBuilder.WriteString(title)
		contentBuilder.WriteString("\n\n")
	}

	// 2. 메타 디스크립션 추출 (우선순위: og:description > twitter:description > description)
	description := ""
	doc.Find("meta[property='og:description']").Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists && description == "" {
			description = content
		}
	})
	if description == "" {
		doc.Find("meta[name='twitter:description']").Each(func(i int, s *goquery.Selection) {
			if content, exists := s.Attr("content"); exists {
				description = content
			}
		})
	}
	if description == "" {
		doc.Find("meta[name='description']").Each(func(i int, s *goquery.Selection) {
			if content, exists := s.Attr("content"); exists {
				description = content
			}
		})
	}
	description = strings.TrimSpace(description)
	if description != "" {
		contentBuilder.WriteString("설명: ")
		contentBuilder.WriteString(description)
		contentBuilder.WriteString("\n\n")
	}

	// 3. 키워드 메타 태그 추출
	keywords := ""
	doc.Find("meta[name='keywords']").Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists && keywords == "" {
			keywords = content
		}
	})
	if keywords != "" {
		contentBuilder.WriteString("키워드: ")
		contentBuilder.WriteString(keywords)
		contentBuilder.WriteString("\n\n")
	}

	// 4. YouTube 특별 처리: YouTube Data API 우선 사용
	if isYouTube {
		videoID := extractYouTubeVideoID(url)
		if videoID != "" {
			// YouTube Data API 시도
			videoInfo, err := getYouTubeVideoInfo(videoID)
			if err == nil && videoInfo != nil {
				// API 성공: API 데이터 우선 사용
				contentBuilder.WriteString("\n=== YouTube API 정보 ===\n")
				contentBuilder.WriteString("동영상 제목: ")
				contentBuilder.WriteString(videoInfo.Title)
				contentBuilder.WriteString("\n\n")

				if videoInfo.Description != "" {
					contentBuilder.WriteString("동영상 설명: ")
					// 설명이 너무 길면 앞부분만 사용
					desc := videoInfo.Description
					if len(desc) > 1000 {
						desc = desc[:1000] + "..."
					}
					contentBuilder.WriteString(desc)
					contentBuilder.WriteString("\n\n")
				}

				if len(videoInfo.Tags) > 0 {
					contentBuilder.WriteString("동영상 태그: ")
					contentBuilder.WriteString(strings.Join(videoInfo.Tags, ", "))
					contentBuilder.WriteString("\n\n")
				}

				log.Printf("[YouTube API] API 데이터 사용 (videoID=%s, title=%s)", videoID, videoInfo.Title)
			} else {
				// API 실패 또는 키 없음: 크롤링 방식으로 대체
				log.Printf("[YouTube API] API 사용 불가, 크롤링 방식으로 진행 (videoID=%s)", videoID)

				// 기존 크롤링 방식
				videoDesc := ""
				doc.Find("meta[name='description']").Each(func(i int, s *goquery.Selection) {
					if content, exists := s.Attr("content"); exists && content != "" && videoDesc == "" {
						videoDesc = content
					}
				})
				if videoDesc != "" && videoDesc != description {
					contentBuilder.WriteString("동영상 상세 설명: ")
					contentBuilder.WriteString(videoDesc)
					contentBuilder.WriteString("\n\n")
				}

				videoKeywords := ""
				doc.Find("meta[name='keywords']").Each(func(i int, s *goquery.Selection) {
					if content, exists := s.Attr("content"); exists && content != "" && videoKeywords == "" {
						videoKeywords = content
					}
				})
				if videoKeywords != "" && videoKeywords != keywords {
					contentBuilder.WriteString("동영상 키워드: ")
					contentBuilder.WriteString(videoKeywords)
					contentBuilder.WriteString("\n\n")
				}
			}
		}

		log.Printf("[AI 태그 추천] YouTube 특별 처리 완료 (url=%s, 현재 길이=%d)", url, contentBuilder.Len())
	}

	// 5. Naver 특별 처리: 동적 페이지 대응
	isNaver := strings.Contains(url, "naver.com") || strings.Contains(url, "navercorp.com")
	if isNaver {
		// Naver 채용 같은 동적 페이지는 메타 태그가 거의 없을 수 있음
		// URL에서 컨텍스트 정보 추출
		if strings.Contains(url, "recruit") || strings.Contains(url, "career") {
			contentBuilder.WriteString("\n페이지 유형: 채용 공고\n")
		}

		// URL 파라미터에서 정보 추출
		if strings.Contains(url, "annoId=") {
			contentBuilder.WriteString("공고 ID가 포함된 URL입니다.\n")
		}

		log.Printf("[AI 태그 추천] Naver 특별 처리 완료 (url=%s, 현재 길이=%d)", url, contentBuilder.Len())
	}

	// 6. 본문 내용 추출 (article, main, section, div 태그)
	initialLen := contentBuilder.Len()
	contentBuilder.WriteString("\n본문 내용:\n")

	// 더 다양한 선택자로 본문 추출 시도
	doc.Find("article, main, [role='main'], .content, .post-content, .article-content, #content").Each(func(i int, s *goquery.Selection) {
		// 광고, 댓글, 네비게이션, 헤더, 푸터 영역 제외
		s.Find("nav, aside, header, footer, .ad, .advertisement, .comment, .sidebar, .menu, script, style").Remove()
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 50 {
			contentBuilder.WriteString(text)
			contentBuilder.WriteString("\n")
		}
	})

	// 본문이 충분하지 않으면 section 태그 시도
	if contentBuilder.Len()-initialLen < 200 {
		doc.Find("section").Each(func(i int, s *goquery.Selection) {
			s.Find("nav, aside, footer, .ad, .advertisement, .comment").Remove()
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 50 {
				contentBuilder.WriteString(text)
				contentBuilder.WriteString("\n")
			}
		})
	}

	// 여전히 부족하면 p 태그에서 직접 추출
	if contentBuilder.Len()-initialLen < 200 {
		doc.Find("p").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 30 {
				contentBuilder.WriteString(text)
				contentBuilder.WriteString("\n")
			}
		})
	}

	// 7. 해시태그 추출
	hashtags := []string{}
	doc.Find("a[href*='tag'], .hashtag, [class*='tag'], [class*='keyword']").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.HasPrefix(text, "#") || (len(text) > 2 && len(text) < 20) {
			hashtags = append(hashtags, text)
		}
	})
	if len(hashtags) > 0 {
		contentBuilder.WriteString("\n해시태그/키워드: ")
		contentBuilder.WriteString(strings.Join(hashtags, ", "))
	}

	content := contentBuilder.String()

	// 메타 태그만으로 충분한 내용이 있으면 진행 (제목 + 설명이 충분한 경우)
	metaContentLength := len(title) + len(description) + len(keywords)

	// 최소한의 정보가 있는지 확인
	// 1. 메타 태그가 있거나
	// 2. 특별 처리된 페이지(YouTube, Naver)에서 컨텍스트가 추가되었거나
	// 3. 전체 콘텐츠가 50자 이상이면 진행
	hasMetaTags := len(title) > 0 || len(description) > 0 || len(keywords) > 0
	hasContext := (isYouTube || isNaver) && len(content) > 50
	hasMinimalInfo := hasMetaTags || hasContext

	// 본문이 너무 짧은 경우 체크
	if len(content) < 100 {
		if hasMinimalInfo {
			log.Printf("[AI 태그 추천] 본문이 부족하지만 메타 태그/컨텍스트로 진행 (url=%s, content_length=%d, meta_length=%d, title_length=%d, desc_length=%d, hasContext=%v)",
				url, len(content), metaContentLength, len(title), len(description), hasContext)
			// 메타 태그 또는 컨텍스트 내용으로 진행 - AI가 판단하도록 함
		} else {
			log.Printf("[AI 태그 추천] 본문 추출 실패: 정보 없음 (url=%s, length=%d, meta_length=%d)", url, len(content), metaContentLength)
			return nil, fmt.Errorf("insufficient content extracted from URL")
		}
	}

	// 본문이 너무 길면 최대 8000자로 제한 (OpenAI 토큰 제한 고려)
	if len(content) > 8000 {
		content = content[:8000] + "..."
	}

	// OpenAI API 호출하여 태그 추천 받기
	tags, err := callOpenAIForTags(content, url)
	if err != nil {
		log.Printf("[AI 태그 추천] OpenAI 태그 생성 실패 (url=%s): %v", url, err)
		return nil, fmt.Errorf("failed to get AI recommendations: %v", err)
	}

	return &AITagRecommendationRes{
		URL:             url,
		RecommendedTags: tags,
	}, nil
}

// callOpenAIForTags OpenAI API를 호출하여 태그 추천을 받습니다
func callOpenAIForTags(content string, url string) ([]string, error) {
	apiKey := localConfig.GetEnvConfig("openaiApiKey")
	if apiKey == "" {
		log.Printf("[AI 태그 추천] OpenAI API 키가 설정되지 않았습니다.")
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	client := openai.NewClient(apiKey)

	// AI 태그 정책에 따른 시스템 프롬프트
	systemPrompt := `당신은 웹 콘텐츠를 분석하여 검색에 유용한 태그를 추천하는 전문가입니다.

다음 규칙을 엄격히 따라 태그를 생성하세요:

**포함해야 할 태그:**
- 핵심 개념/명사형 주제어 (예: 프롬프트, UX리서치, 포트폴리오)
- 고유명사 (예: ChatGPT, Notion, 토스, Apple, Figma)
- 전문 용어 (예: SBI모델, 파인튜닝, 데이터레이블링)
- 도메인 단어 (예: 클라우드, 마케팅, 디자인) - 문맥상 의미 있을 때만
- 본문 내 해시태그

**절대 제외해야 할 태그:**
- 도메인명/플랫폼명 (예: yozm, naver, brunch, medium)
- 감정 형용사 (예: 멋진, 새로운, 유용한)
- 문장형 표현 (예: ~하는 방법, ~를 해보자)
- CTA 문구 (예: 클릭, 확인, 공유, 신청하기)
- 날짜/버전 정보 (예: 2024, 1.0, 7월)
- 이모지/특수문자

**형식 규칙:**
- 명사형 중심으로 작성
- 1~10자 사이의 짧고 직관적인 단어
- 최대 5개까지만 추천
- JSON 배열 형식으로만 응답: ["태그1", "태그2", ...]

콘텐츠의 핵심 주제를 파악하고, 사용자가 나중에 검색할 때 유용한 태그만 선정하세요.`

	userPrompt := fmt.Sprintf("URL: %s\n\n콘텐츠:\n%s\n\n위 콘텐츠를 분석하여 최대 5개의 추천 태그를 JSON 배열로 반환하세요.", url, content)

	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: "gpt-4o-mini", // GPT-4o mini 모델 사용
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Temperature:           0.3, // 일관성을 위해 낮은 temperature 사용
		MaxCompletionTokens:   200,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("[AI 태그 추천] OpenAI API 호출 실패 (url=%s): %v", url, err)
		return nil, fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Printf("[AI 태그 추천] OpenAI 응답이 비어 있습니다 (url=%s)", url)
		return nil, fmt.Errorf("no response from OpenAI")
	}

	responseText := strings.TrimSpace(resp.Choices[0].Message.Content)

	// 마크다운 코드 블록 제거 (```json ... ``` 형식)
	if strings.HasPrefix(responseText, "```") {
		// 첫 번째 줄바꿈 이후부터 마지막 ``` 이전까지 추출
		lines := strings.Split(responseText, "\n")
		if len(lines) > 2 {
			// 첫 줄(```json)과 마지막 줄(```) 제거
			responseText = strings.Join(lines[1:len(lines)-1], "\n")
			responseText = strings.TrimSpace(responseText)
			log.Printf("[AI 태그 추천] 마크다운 코드 블록 제거 완료 (url=%s)", url)
		}
	}

	// JSON 배열 파싱
	var tags []string
	err = json.Unmarshal([]byte(responseText), &tags)
	if err != nil {
		// JSON 파싱 실패 시 응답에서 태그 추출 시도
		// 예: "["태그1", "태그2"]" 형식이 아닌 경우 처리
		log.Printf("[AI 태그 추천] OpenAI 응답 파싱 실패 (url=%s, response=%s, err=%v)", url, responseText, err)
		return nil, fmt.Errorf("failed to parse AI response as JSON: %v, response: %s", err, responseText)
	}

	// 최대 5개로 제한
	if len(tags) > 5 {
		tags = tags[:5]
	}

	return tags, nil
}
