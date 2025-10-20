package link

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/errgo.v2/errors"
	localConfig "joosum-backend/pkg/config"
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

func (u LinkUsecase) FindAllLinksByUserIdAndLinkBookId(userId string, linkBookId string, sort string, order string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserIdAndLinkBookId(userId, linkBookId, sort, order)
	if err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return []*Link{}, nil
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

// GetAIRecommendedTags AI를 사용하여 URL의 본문 내용을 분석하고 추천 태그를 생성합니다
func (LinkUsecase) GetAIRecommendedTags(url string) (*AITagRecommendationRes, error) {
	// URL 이 http:// 혹은 https:// 로 시작하지 않으면 https:// 를 붙입니다.
	url = util.EnsureHTTPPrefix(url)

	// URL에서 본문 내용 크롤링
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// 본문 텍스트 추출
	var contentBuilder strings.Builder

	// 1. 제목 추출 (og:title 또는 title 태그)
	title := ""
	doc.Find("meta[property='og:title']").Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists {
			title = content
		}
	})
	if title == "" {
		title = doc.Find("title").First().Text()
	}
	contentBuilder.WriteString("제목: ")
	contentBuilder.WriteString(title)
	contentBuilder.WriteString("\n\n")

	// 2. 메타 디스크립션 추출
	doc.Find("meta[property='og:description'], meta[name='description']").Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists && i == 0 {
			contentBuilder.WriteString("설명: ")
			contentBuilder.WriteString(content)
			contentBuilder.WriteString("\n\n")
		}
	})

	// 3. 본문 내용 추출 (article, main, section, p 태그)
	contentBuilder.WriteString("본문 내용:\n")
	doc.Find("article, main, section").Each(func(i int, s *goquery.Selection) {
		// 광고, 댓글, 네비게이션 영역 제외
		s.Find("nav, aside, footer, .ad, .advertisement, .comment").Remove()
		text := strings.TrimSpace(s.Text())
		if text != "" && len(text) > 50 {
			contentBuilder.WriteString(text)
			contentBuilder.WriteString("\n")
		}
	})

	// section이 없는 경우 p 태그에서 직접 추출
	if contentBuilder.Len() < 200 {
		doc.Find("p").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 30 {
				contentBuilder.WriteString(text)
				contentBuilder.WriteString("\n")
			}
		})
	}

	// 4. 해시태그 추출
	hashtags := []string{}
	doc.Find("a[href*='tag'], .hashtag, [class*='tag']").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.HasPrefix(text, "#") {
			hashtags = append(hashtags, text)
		}
	})
	if len(hashtags) > 0 {
		contentBuilder.WriteString("\n해시태그: ")
		contentBuilder.WriteString(strings.Join(hashtags, ", "))
	}

	content := contentBuilder.String()

	// 본문이 너무 짧으면 에러
	if len(content) < 100 {
		return nil, fmt.Errorf("insufficient content extracted from URL")
	}

	// 본문이 너무 길면 최대 8000자로 제한 (OpenAI 토큰 제한 고려)
	if len(content) > 8000 {
		content = content[:8000] + "..."
	}

	// OpenAI API 호출하여 태그 추천 받기
	tags, err := callOpenAIForTags(content, url)
	if err != nil {
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
		Model: "gpt-5-nano", // GPT-5 nano 모델 사용
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
		Temperature: 0.3, // 일관성을 위해 낮은 temperature 사용
		MaxTokens:   200,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	responseText := strings.TrimSpace(resp.Choices[0].Message.Content)

	// JSON 배열 파싱
	var tags []string
	err = json.Unmarshal([]byte(responseText), &tags)
	if err != nil {
		// JSON 파싱 실패 시 응답에서 태그 추출 시도
		// 예: "["태그1", "태그2"]" 형식이 아닌 경우 처리
		return nil, fmt.Errorf("failed to parse AI response as JSON: %v, response: %s", err, responseText)
	}

	// 최대 5개로 제한
	if len(tags) > 5 {
		tags = tags[:5]
	}

	return tags, nil
}
