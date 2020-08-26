IMAGE ?= jasonmccallister/httputil
TAG ?= latest

local:
	docker-compose down -v
	docker-compose up -d --build
	docker-compose ps
deploy: build push

build:
	docker build -t ${IMAGE}:${TAG} .
push:
	docker push ${IMAGE}:${TAG}
