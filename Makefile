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

VERSION=0.0.5
IMAGE=imkira/gcp-iap-auth

.PHONY: all build build-docker docker-release dist clean deps vendor

all: build

build:
	go build -o "build/gcp-iap-auth"

docker-build:
	docker build -t "${IMAGE}:${VERSION}" .

docker-login:
	echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

docker-release:
	docker tag "${IMAGE}:${VERSION}" "${IMAGE}:latest"
	docker push "${IMAGE}:${VERSION}"
	docker push "${IMAGE}:latest"

dist:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD_ENV) go build $(GOBUILD_FLAGS) -o "dist/$(DIST_OUTPUT)"
	cd dist && $(SHASUM) "$(DIST_OUTPUT)" > "$(DIST_OUTPUT).sha1"

clean:
	rm -rf build dist
