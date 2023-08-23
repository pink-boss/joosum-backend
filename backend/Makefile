# config.yml 에서 변수를 가져옴
IMAGE_SERVER := $(shell sed -n 's/^imageServer: //p' config.yml)
IMAGE_NAME := joosum
CONTAINER_NAME := joosum

docker.dev:
	swag init
	docker build -t ${IMAGE_NAME} .
	docker run --rm -d -v /etc/ssl/certs:/etc/ssl/certs -p 5001:5001 \
		--name ${CONTAINER_NAME}_dev ${IMAGE_NAME}

docker.prod:
	swag init
	docker build -t ${IMAGE_NAME} .
	docker run --rm -d -v /etc/ssl/certs:/etc/ssl/certs -p 5001:5001 \
		--log-opt max-size=3k --log-opt max-file=1 \
 		--name ${CONTAINER_NAME}_prod ${IMAGE_NAME} -env=prod

docker.stop:
	-docker stop ${CONTAINER_NAME}_dev # 첫번째 라인이 실패해도 두번째라인을 실행
	docker stop ${CONTAINER_NAME}_prod

# 수동배포를 위해 이미지를 푸시합니다. ex) make docker.push v=230514-2 혹은 make docker.push v=1.3.2
docker.push:
	docker build -t "$(IMAGE_SERVER):$(v)" .
	docker push $(IMAGE_SERVER):$(v)
