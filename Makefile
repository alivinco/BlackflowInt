GOPATH=$(shell pwd)/../../../../
export GOPATH
export PATH=$PATH:$GOPATH/bin

.PHONY : clean

clean :
	rm -Rf $(GOPATH)bin/*
	rm -Rf bin

install-glide :
	curl https://glide.sh/get | sh

init-project :
	glide init

build :
	go install


docker-clean:
	docker rm blacktower
	docker rmi alivinco/blacktower

.PHONY : dist-docker

dist-docker :
	mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 /usr/local/go/bin/go build -v -o bin/blackflowint
	echo $(shell ls -a bin/)
	docker build -t alivinco/blackflowint .

docker-publish : dist-docker
	docker push alivinco/blackflowint
