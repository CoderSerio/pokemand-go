VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -X github.com/CoderSerio/pokemand-go/cmd.Version=$(VERSION)
PLATFORMS := linux/amd64 darwin/amd64 windows/amd64
BINARY_NAME := pkmg

.PHONY: all build clean install release

all: build

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)

install:
	go build -ldflags "$(LDFLAGS)" -o $(GOPATH)/bin/$(BINARY_NAME)

clean:
	rm -rf bin/
	rm -rf dist/


# 发布新版本
release:
	goreleaser release --clean


# 发布多平台二进制文件
# release: clean
# 	mkdir -p dist
# 	$(foreach platform,$(PLATFORMS),\
# 		GOOS=$(word 1,$(subst /, ,$(platform))) \
# 		GOARCH=$(word 2,$(subst /, ,$(platform))) \
# 		go build -ldflags "$(LDFLAGS)" \
# 		-o dist/$(BINARY_NAME)-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform)))$(if $(findstring windows,$(platform)),.exe,) ;\
# 	) 