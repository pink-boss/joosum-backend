.PHONY: clean test security build run swag

APP_NAME = apiserver
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

docker.build:
	docker build -t net_http-go-template .

docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network

docker.run: docker.network swag docker.net_http

# 서버 컨테이너를 내렸다가 다시 띄움
docker.restart: docker.stop docker.run

docker.stop:
	docker stop cgapp-net_http

docker.net_http: docker.build
	docker run --rm -d \
		--name cgapp-net_http \
		--network dev-network \
		-p 5001:5001 \
		net_http-go-template

swag:
	swag init
