GOARCH=amd64
LDFLAGS := -s -w
PROJECT_NAME=nat-tcp-server

DOCKER_URL=registry.cn-hangzhou.aliyuncs.com/xxcheng
DOCKER_USERNAME=developer@xxcheng.cn

VERSION=$(shell git describe --tags --always)

.PHONY: build-server-linux
build-linux: # Build project for Linux | 构建Linux下的可执行文件
	env CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o ./build/$(PROJECT_NAME) ./server/main.go
	@echo "Build project for Linux successfully"

.PHONY: pubulish-server-docker
publish-docker:
	docker build -t $(PROJECT_NAME):$(VERSION) -f Dockerfile-server .;
	docker login --username=$(DOCKER_USERNAME) $(DOCKER_URL);
	docker tag $(PROJECT_NAME):$(VERSION) $(DOCKER_URL)/$(PROJECT_NAME):$(VERSION);
	docker tag $(PROJECT_NAME):$(VERSION) $(DOCKER_URL)/$(PROJECT_NAME):latest;
	docker push $(DOCKER_URL)/$(PROJECT_NAME):$(VERSION);
	docker push $(DOCKER_URL)/$(PROJECT_NAME):latest;
	@echo "Publish docker successfully"

.PHONY: fast-server-docker
fast-docker:
	make build-linux;
	make publish-docker;