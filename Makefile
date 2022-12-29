REGISTRY := hub.docker.com/repository/docker/

REPOSITORY := futtie/legocollector
PWD := $(shell pwd)
#TAG := $(shell git symbolic-ref --short HEAD | egrep -o '((PBX|DRIFT)-([0-9]+))|(^master)' || echo 'UnknownBranch')
TAG := 0.1

build: 
	docker build --no-cache -t ${REPOSITORY}:local -f Dockerfile .

build_local: 
	docker build --no-cache -f Dockerfile .

#build: docker_login
#	docker build --no-cache 
#	  -t ${REGISTRY}/${REPOSITORY}:${TAG} \
#	  -t ${REPOSITORY}:local -f Dockerfile \
#	  .

lint:
	docker run --rm -v "${PWD}":/app golang:1.13 sh -c "cd /app && go get golang.org/x/lint/golint && golint -set_exit_status ./..."

push: docker_login
	docker push ${REGISTRY}/${REPOSITORY}:${TAG} 
	docker logout ${REGISTRY}

run: build
	docker-compose rm -f
	docker-compose up 

docker_login:
	docker login -u futtie -p On1gU*iX32At@snYR*YcHc%o0gd*@5b6 ${REGISTRY}


