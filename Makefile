#
#   这一个非常文艺的编译程序
#

APP_NAME = altar


#一些变量参数
BUILD_PATH := $(shell pwd)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_AUTHOR := $(shell echo "`git config user.name` <`git config user.email`>")
BUILD_TYPE := make

LDFLAGS := -X 'main.BuildPath=$(BUILD_PATH)' \
           -X 'main.BuildTime=$(BUILD_TIME)' \
           -X 'main.BuildAuthor=$(BUILD_AUTHOR)' \
           -X 'main.BuildType=$(BUILD_TYPE)'

all: mod altar

mod:
	GOPROXY=https://goproxy.cn go mod tidy

altar:
	go build -o altar -ldflags "$(LDFLAGS)" *.go

clean:
	rm altar


.PHONY: all mod altar clean

