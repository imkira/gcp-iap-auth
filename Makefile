GIT_BRANCHTAG=$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
GIT_COMMIT=$(shell git rev-parse HEAD)
GOOS?=linux
GOARCH?=amd64

GOBUILD_ENV+=CGO_ENABLED=0
GOBUILD_FLAGS+=-a -ldflags "$(GOBUILD_LDFLAGS)"
GOBUILD_LDFLAGS+=-X \"main.version=$(GIT_BRANCHTAG)\"
GOBUILD_LDFLAGS+=-X \"main.revision=$(GIT_COMMIT)\"
GOBUILD_LDFLAGS+=-s -w

DIST_OUTPUT=gcp-iap-auth-$(GOOS)-$(GOARCH)

SHASUM=shasum -a 1

.PHONY: all build dist clean deps vendor

all: build

build:
	go build -o "build/gcp-iap-auth"

dist:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD_ENV) go build $(GOBUILD_FLAGS) -o "dist/$(DIST_OUTPUT)"
	$(SHASUM) "dist/$(DIST_OUTPUT)" > "dist/$(DIST_OUTPUT).sha1"

clean:
	rm -rf build dist

deps:
	if [ ! -d "${GOPATH}/src/github.com/Masterminds/glide" ]; then go get -u github.com/Masterminds/glide; fi

vendor:
	glide up -v
