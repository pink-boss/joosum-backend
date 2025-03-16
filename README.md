# :link: JOOSUM

## 실행

---

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
-
