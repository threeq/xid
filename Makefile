# These are the values we want to pass for VERSION  and BUILD
VERSION=v2.0.0
BUILD=`date +%FT%T%z`
# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION} -X main.Build=${BUILD}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build $(LDFLAGS)
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOTOOL=$(GOCMD) tool

BIN_PATH=bin
BINARY_NAME=xid
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_MAC=$(BINARY_NAME)_mac
BINARY_WIN=$(BINARY_NAME)_win.exe
OUTPUT=out
SRC=github.com/threeq/xid/server

.PHONY: all
all: clean test build-all

test:
	$(GOTEST) -v ./...

benchmark:
	$(GOTEST) -v ./...
	@echo "运行 server 性能测试"
	$(GOTEST) -timeout 2h -bench BenchmarkAll -benchmem -cpuprofile cpu.out -memprofile mem.out -run=^$$ internal/rserver/*
	$(GOTOOL) pprof -svg ./rserver.test cpu.out > cpu.benchmarkall.svg
	$(GOTOOl) pprof -svg ./rserver.test mem.out > mem.benchmarkall.svg

clean:
	$(GOCLEAN)
	rm -rf $(OUTPUT)
	rm -rf $(BIN)

build-linux:
	GOPROXY=https://goproxy.io,direct CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY_LINUX) -v $(SRC)
	upx -9 $(BIN_PATH)/$(BINARY_LINUX)
build-mac:
	GOPROXY=https://goproxy.io,direct CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY_MAC)   -v $(SRC)
	upx -9 $(BIN_PATH)/$(BINARY_MAC)
build-win:
	GOPROXY=https://goproxy.io,direct CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY_WIN)   -v $(SRC)

build-all: clean build-linux build-mac build-win
	echo 'build completed.'

# docker-build:
#	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/$(SRC) golang:latest go build -o "$(BINARY_UNIX)" -v

tag: clean build-linux
	git add .
	git commit -am "deploy version $(VERSION)"
	git push origin master
	git tag -a $(VERSION) -m"deploy version $(VERSION)"
	git push --tags

push:
	docker build . -t threewq/xid:latest
	docker tag threewq/xid:latest threewq/xid:$(VERSION)
	docker push threewq/xid:latest
	docker push threewq/xid:$(VERSION)

deploy: tag push
	echo 'deploy completed'

stat: cloc gocyclo
	@echo "代码行数统计"
	@ls *.go main/* Makefile | xargs cloc --by-file
	@echo "圈复杂度统计"
	@ls *.go main/*.go | grep -v _test | xargs gocyclo
	@ls *.go main/*.go | grep -v _test | xargs gocyclo | awk '{sum+=$$1}END{printf("总圈复杂度: %s", sum)}'

cloc:
	@hash cloc 2>/dev/null || { \
        echo "安装代码统计工具 cloc" && \
        mkdir -p third && cd third && \
        wget https://github.com/AlDanial/cloc/archive/v1.82.zip && \
        unzip v1.76.zip; \
    }


gocyclo:
	@hash gocyclo 2>/dev/null || { \
        echo "安装代码圈复杂度统计工具 gocyclo" && \
        go get -u github.com/fzipp/gocyclo; \
    }