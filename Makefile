IMAGE ?= jasonmccallister/httputil
TAG ?= latest

deploy: build tag push

build:
	docker build -t ${IMAGE}:${TAG} .
tag:

push:
	docker push ${IMAGE}:${TAG}
