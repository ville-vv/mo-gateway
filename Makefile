DOCKER_BUILD=./docker_build
BINARY_DIR=./bin
PROJECTNAME=$(shell basename "$(PWD)")
DOCKER_BINARY=$(DOCKER_BUILD)/$(PROJECTNAME)
UNIX_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_linux
MAC_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_mac
WIN_BINARY=$(BINARY_DIR)/$(PROJECTNAME)_amd64_win.exe

.PHONY:all test image clean build 

all: build

test:
	go test  -cover ./...

build: clean build-unix build-osx build-win

build-unix:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(UNIX_BINARY)
build-osx:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(MAC_BINARY)
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(WIN_BINARY)

docker-read:
	sudo rm -rf $(DOCKER_BUILD)
	mkdir -p $(DOCKER_BUILD)
	cp -r ./conf  $(DOCKER_BUILD)
	cp docker/* $(DOCKER_BUILD)
	cp -r $(UNIX_BINARY) $(DOCKER_BUILD)/$(PROJECTNAME)
	echo $(PROJECTNAME)

docker-build: docker-read
	sudo docker build -t $(PROJECTNAME) $(DOCKER_BUILD)

docker_stop:
	sudo docker stop $(sudo docker container ls -a | grep "tstl" | awk '{print $1}')
docker_rm:
	sudo docker rm $(sudo docker container ls -a | grep "tstl" | awk '{print $1}')
docker_rmi:
	sudo docker rmi $(sudo docker images | grep "none" | awk '{print $3}')

clean:
	rm -rf bin
	rm -rf docker_build
