VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -X github.com/CoderSerio/pokemand-go/cmd.Version=$(VERSION)
PLATFORMS := linux/amd64 darwin/amd64 windows/amd64
BINARY_NAME := pkmg

.PHONY: all build clean install release release-dry-run publish-go publish

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

release-dry-run:
	@test -n "$(RELEASE_VERSION)" || (echo "Usage: make release-dry-run RELEASE_VERSION=v0.1.0" && exit 1)
	@git rev-parse "$(RELEASE_VERSION)" >/dev/null 2>&1 && (echo "Tag already exists: $(RELEASE_VERSION)" && exit 1) || true
	@set -e; \
	git tag "$(RELEASE_VERSION)"; \
	trap 'git tag -d "$(RELEASE_VERSION)" >/dev/null 2>&1 || true' EXIT; \
	goreleaser release --clean --skip=publish,announce

publish-go:
	@test -n "$(RELEASE_VERSION)" || (echo "Usage: make publish-go RELEASE_VERSION=v0.2.1" && exit 1)
	@./scripts/publish/go.sh "$(RELEASE_VERSION)"

publish:
	@test -n "$(RELEASE_VERSION)" || (echo "Usage: make publish RELEASE_VERSION=v0.2.1" && exit 1)
	@./scripts/publish/index.sh "$(RELEASE_VERSION)"


# 发布多平台二进制文件
# release: clean
# 	mkdir -p dist
# 	$(foreach platform,$(PLATFORMS),\
# 		GOOS=$(word 1,$(subst /, ,$(platform))) \
# 		GOARCH=$(word 2,$(subst /, ,$(platform))) \
# 		go build -ldflags "$(LDFLAGS)" \
# 		-o dist/$(BINARY_NAME)-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform)))$(if $(findstring windows,$(platform)),.exe,) ;\
# 	) 
