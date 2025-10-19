# AI 태그 추천 기능

## 개요

URL의 본문 내용을 분석하여 검색에 유용한 태그를 AI가 자동으로 최대 5개까지 추천하는 기능입니다.

## 설정

### OpenAI API 키 설정

`backend/config.yml` 파일에 OpenAI API 키를 추가하세요:

```yaml
openaiApiKey: sk-your-openai-api-key-here
```

## API 엔드포인트

### POST /links/ai-tags

URL의 본문을 분석하여 AI 태그를 추천받습니다.

**인증**: JWT Bearer Token 필요

**요청 본문**:
```json
{
  "url": "https://brunch.co.kr/@wine-ny/163"
}
```

**응답**:
```json
{
  "url": "https://brunch.co.kr/@wine-ny/163",
  "recommendedTags": [
    "프롬프트",
    "요구사항분석",
    "기획자",
    "개발자협업",
    "기능명세서"
  ]
}
```

**에러 응답**:
- 400: 잘못된 요청 본문
- 401: Authorization 헤더 없음
- 500: AI 태그 추천 과정에서 오류 발생

## 태그 생성 정책

AI는 `ai_tag_policy.md`에 정의된 정책에 따라 태그를 생성합니다:

### ✅ 포함되는 태그
- 핵심 개념/명사형 주제어 (예: 프롬프트, UX리서치, 포트폴리오)
- 고유명사 (예: ChatGPT, Notion, 토스, Apple, Figma)
- 전문 용어 (예: SBI모델, 파인튜닝, 데이터레이블링)
- 도메인 단어 (예: 클라우드, 마케팅, 디자인)
- 본문 내 해시태그

### ❌ 제외되는 태그
- 도메인명/플랫폼명 (예: yozm, naver, brunch, medium)
- 감정 형용사 (예: 멋진, 새로운, 유용한)
- 문장형 표현 (예: ~하는 방법, ~를 해보자)
- CTA 문구 (예: 클릭, 확인, 공유)
- 날짜/버전 정보 (예: 2024, 1.0)
- 이모지/특수문자

## 사용 방법

### 1. Swagger UI에서 테스트

1. 서버 실행: `go run main.go`
2. 브라우저에서 `http://localhost:5001/swagger/index.html` 접속
3. `/auth/google` 또는 `/auth/signup`으로 로그인하여 JWT 토큰 획득
4. 우측 상단 "Authorize" 버튼 클릭
5. `Bearer {your-access-token}` 형식으로 토큰 입력
6. `/links/ai-tags` POST 엔드포인트에서 "Try it out" 클릭
7. URL 입력 후 "Execute" 클릭

### 2. cURL로 테스트

```bash
curl -X POST "http://localhost:5001/links/ai-tags" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://brunch.co.kr/@wine-ny/163"
  }'
```

### 3. 테스트 예시 URL

다음 URL들로 테스트해볼 수 있습니다:

- **기획/개발**: https://brunch.co.kr/@wine-ny/163
- **디자인**: https://yozm.wishket.com/magazine/detail/2700/
- **마케팅**: https://blog.toss.im/article/toss-brand-marketing
- **기술블로그**: https://engineering.linecorp.com/ko/blog/

## 작동 원리

1. **URL 크롤링**: URL에서 HTML을 가져와 파싱합니다
2. **본문 추출**:
   - 제목 (og:title, title 태그)
   - 메타 디스크립션 (og:description, description)
   - 본문 내용 (article, main, section, p 태그)
   - 해시태그 (#로 시작하는 링크)
3. **광고/댓글 제거**: nav, aside, footer, .ad, .advertisement, .comment 영역 제외
4. **OpenAI API 호출**: GPT-4o-mini 모델로 태그 추천 요청
5. **태그 반환**: JSON 배열 형식으로 최대 5개 태그 반환

## 주의사항

1. OpenAI API 키가 설정되지 않으면 500 에러가 발생합니다
2. URL이 접근 불가하거나 본문이 너무 짧으면 에러가 발생합니다
3. 본문이 4000자를 초과하면 잘라서 분석합니다 (토큰 제한)
4. 크롤링이 금지된 사이트는 분석할 수 없습니다

## 파일 구조

```
backend/
├── app/link/
│   ├── link_usecase.go        # GetAIRecommendedTags(), callOpenAIForTags()
│   ├── link_handler.go        # GetAIRecommendedTags handler
│   └── link_model.go          # AITagRecommendationReq, AITagRecommendationRes
├── pkg/routes/
│   └── private_routes.go      # POST /links/ai-tags 라우트
├── config.yml                 # openaiApiKey 설정
└── ai_tag_policy.md          # AI 태그 생성 정책 문서
```

## OpenAI 모델

- **모델**: GPT-4o-mini
- **Temperature**: 0.3 (일관성을 위해 낮게 설정)
- **MaxTokens**: 200
- **응답 형식**: JSON 배열 `["태그1", "태그2", ...]`

## 에러 처리

에러 발생 시 로그에 상세한 정보가 출력됩니다:

```json
{
  "timestamp": "...",
  "level": "error",
  "msg": "Response with errors",
  "status_code": 500,
  "response_body": {"code": 1002, "message": "서버 오류가 발생했습니다."},
  "errors": "Error #01: GetAIRecommendedTags failed: ..."
}
```
