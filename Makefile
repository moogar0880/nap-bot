BIN_NAME=nap-bot

VERSION=$(shell git describe --tags 2> /dev/null || echo '0.0.0')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
PRERELEASE=
IMAGE_NAME := moogar0880/nap-bot

# if we have untagged commits, mark this build as a pre-release
ifneq ($(strip $(GIT_DIRTY)),)
PRERELEASE=DEV
endif
.PHONY: all
all: clean vendor binary test


.PHONY: help
help:
	@echo 'Management commands for nap-bot:'
	@echo
	@echo 'Usage:'
	@echo '    make clean           Clean the directory tree.'
	@echo '    make test            Run tests on the project.'
	@echo '    make test/benchmark  Run benchmark tests on the project.'
	@echo '    make test/coverage   Generate and view HTML coverage report in a browser.'
	@echo '    make vendor          runs dep to fetch vendor dependencies.'
	@echo '    make binary          Compile the binary for this project.'
	@echo '    make build/alpine    Compile optimized for alpine linux.'
	@echo '    make package         Build final docker image with just the go binary inside'
	@echo '    make push            Push tagged images to registry'
	@echo '    make tag             Tag image created by package with latest, git commit and version'
	@echo

##############################################################################
# The following targets are used for aiding in development and CI for the 
# nap-bot source code
##############################################################################
.PHONY: clean
clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME} && go clean


.PHONY: test
test:
	go test -coverprofile=coverage.out ./...

.PHONY: test/benchmark
test/benchmark:
	go test -run=XXX -bench=. -benchmem

.PHONY: test/coverage
test/coverage: test
	go tool cover -html=coverage.out

.PHONY: vendor
vendor:
	dep ensure


##############################################################################
# The following targets are used for packaging the nap-bot
# binary into a docker container
##############################################################################
.PHONY: binary
binary:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.VersionPrerelease=${PRERELEASE}" -o bin/${BIN_NAME}

.PHONY: build/alpine
build/alpine:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	CGO_ENABLED=0 go build \
		-ldflags '-w -linkmode external -extldflags "-static" -X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.VersionPrerelease=VersionPrerelease=RC' \
		-o bin/${BIN_NAME}

.PHONY: package
package:
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(IMAGE_NAME):local .

.PHONY: tag
tag: 
	@echo "Tagging: latest ${VERSION} $(GIT_COMMIT)"
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):$(GIT_COMMIT)
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):${VERSION}
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):latest

.PHONY: push
push: tag
	@echo "Pushing docker image to registry: latest ${VERSION} $(GIT_COMMIT)"
	docker push $(IMAGE_NAME):$(GIT_COMMIT)
	docker push $(IMAGE_NAME):${VERSION}
	docker push $(IMAGE_NAME):latest
	gcloud docker -- gcr.io/${IMAGE_NAME}:${VERSION}

.PHONY: package/gcloud
package/gcloud:
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) -t gcr.io/$(IMAGE_NAME):local .

.PHONY: tag/gcloud
tag/gcloud: package/gcloud
	@echo "Tagging: latest ${VERSION}"
	docker tag gcr.io/$(IMAGE_NAME):local gcr.io/$(IMAGE_NAME):${VERSION}

.PHONY: push/gcloud
push/gcloud: tag/gcloud
	@echo "Pushing docker image to gcr: ${VERSION}"
	gcloud docker -- push gcr.io/${IMAGE_NAME}:${VERSION}
