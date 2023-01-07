REGISTRY := https://hub.docker.com/repository/docker

REPOSITORY := futtie/legocollector
PWD := $(shell pwd)
TAG := 0.4

build: 
	docker build --no-cache -t ${REPOSITORY}:local -f Dockerfile .

build_local: 
	docker build -f Dockerfile .

buildtag: docker_login
	docker build --no-cache -t ${REPOSITORY}:${TAG} -t ${REPOSITORY}:latest -t ${REPOSITORY}:local -f Dockerfile .

lint:
	docker run --rm -v "${PWD}":/app golang:1.13 sh -c "cd /app && go get golang.org/x/lint/golint && golint -set_exit_status ./..."

push: docker_login
	docker push ${REPOSITORY}:${TAG} 
	docker push ${REPOSITORY}:latest 
	docker logout 

run: build
	docker-compose rm -f
	docker-compose up 

docker_login:
	cat /home/mhe/dockerpassword.txt | docker login -u futtie --password-stdin

