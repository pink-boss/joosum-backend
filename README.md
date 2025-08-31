# :link: JOOSUM

## 실행
d
로컬 실행

```bash
go run main.go # 개발환경으로 실행
go run main.go -env=prod # 운영환경으로 실행
```

<br>
도커 실행

```bash
make docker.dev # 개발환경으로 실행
make docker.prod # 운영환경으로 실행
```

## 로그

---

<br>
로그를 json format 으로 보기위해 jq 설치 (optional)

```bash
brew install jq
```

<br>
로그 확인

```bash
docker logs joosum_dev -f # 개발환경 로그확인
docker logs server 2>&1 -f | jq # 배포서버 로그확인
```

## Version

---

- Mongo-driver: v1.11.4

## 오류 코드 체계

| 코드 | 메시지 |
|------|------------------------|
| 1000 | 잘못된 요청 본문입니다. |
| 1001 | Authorization 헤더가 없습니다. |
| 1002 | 서버 오류가 발생했습니다. |
| 1003 | 필수 파라미터가 누락되었습니다. |
| 2000 | 유효하지 않은 ID 토큰입니다. |
| 2001 | 이미 존재하는 사용자입니다. |
| 2002 | 탈퇴 후 30일이 지나지 않았습니다. |
| 3000 | 같은 이름의 폴더가 존재합니다. |

