ROOT_PATH=$(shell pwd)
RELEASE=/bin/
APP_NAME=stock_stat
GO_LIST=$(shell go list ${ROOT_PATH} | grep -v vendor)
GO_FILE=$(shell find . -name '*.go' | grep -v vendor)

.PHONY : build

all : clean proc

init :
	@go mod tidy

fmt :
	@gofmt -l -w ${GO_FILE}

test :
	@go test -cover ${GO_LIST}

proc :
	cd ${ROOT_PATH} && \
	go build -o ${ROOT_PATH}${RELEASE}${APP_NAME}


clean:
	@rm -rf ${ROOT_PATH}${RELEASE}/*