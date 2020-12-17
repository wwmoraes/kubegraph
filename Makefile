ICONS_FOLDER := icons
ICONS_PKG := icons
ICONS_GO_FILE := $(ICONS_FOLDER)/icons.go
ICONS_FILES := $(ICONS_FOLDER)/*.svg
CMD_SOURCE_FILES := $(shell find cmd -type f -name '*.go')
INTERNAL_SOURCE_FILES := $(shell find internal -type f -name '*.go')
ICONS_SOURCE_FILES := $(wildcard icons/*.go)
SOURCE_FILES := $(CMD_SOURCE_FILES) $(INTERNAL_SOURCE_FILES) $(ICONS_SOURCE_FILES)

GIT_SHA = sha-$(shell git log -n 1 --format="%h")
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_REV = $(shell git log -n 1 --format="%H")
DATE = $(shell date -u +"%Y-%m-%dT%TZ")

REPO := wwmoraes/kubegraph
USERNAME = $(shell git config user.name)
EMAIL = $(shell git config user.email)
OCI_TITLE = kubegraph
OCI_DESCRIPTION = Kubernetes resource graph generator
OCI_URL = https://github.com/$(REPO)
OCI_SOURCE = https://github.com/$(REPO)
OCI_VERSION = $(GIT_BRANCH)
OCI_CREATED = $(DATE)
OCI_REVISION = $(GIT_REV)
OCI_LICENSES = MIT
OCI_AUTHORS = $(USERNAME) <$(EMAIL)>
OCI_DOCUMENTATION = https://github.com/$(REPO)
OCI_AUTHORS = $(USERNAME) <$(EMAIL)>


.DEFAULT_GOAL := build

.PHONY: icons
icons: $(ICONS_GO_FILE)

$(ICONS_GO_FILE): $(ICONS_FILES)
	@type -p go-bindata 1>/dev/null || go get -u github.com/go-bindata/go-bindata/...
	$(info removing old icons package file...)
	-@rm $(ICONS_GO_FILE)
	$(info regenerating icons package file...)
	@go-bindata -o $(ICONS_GO_FILE) -pkg $(ICONS_PKG) -nometadata -mode 0664 -nomemcopy $(ICONS_FOLDER)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -v ./...

.PHONY: coverage
coverage: coverage.out
	@go tool cover -func=$<

.PHONY: coverage-html
coverage-html: coverage.html

%.html: %.out
	go tool cover -html=$< -o $@

%.out: $(SOURCE_FILES)
	@go test -race -cover -coverprofile=$@ -v ./...

build: kubegraph

kubegraph: $(SOURCE_FILES) vendor
	go build -mod=vendor -race -o ./ ./...

vendor: go.mod go.sum
	go mod vendor

.PHONY: run
run:
	go run cmd/kubegraph/main.go sample.yaml
	dot -Tsvg -o sample.svg sample.dot

.PHONY: image
image: Dockerfile $(SOURCE_FILES)
	@docker build \
		--label org.opencontainers.image.title=$(OCI_TITLE) \
		--label org.opencontainers.image.description="$(OCI_DESCRIPTION)" \
		--label org.opencontainers.image.url=$(OCI_URL) \
		--label org.opencontainers.image.source=$(OCI_SOURCE) \
		--label org.opencontainers.image.version=$(OCI_VERSION) \
		--label org.opencontainers.image.created=$(OCI_CREATED) \
		--label org.opencontainers.image.revision=$(OCI_REVISION) \
		--label org.opencontainers.image.licenses=$(OCI_LICENSES) \
		--label org.opencontainers.image.authors="$(OCI_AUTHORS)" \
		--label org.opencontainers.image.documentation=$(OCI_DOCUMENTATION) \
		--label org.opencontainers.image.vendor="$(OCI_VENDOR)" \
		--cache-from $(REPO):single-$(GIT_SHA) \
		--cache-from $(REPO):single-$(GIT_BRANCH) \
		--cache-from $(REPO):single-master \
		--cache-from $(REPO):single-latest \
  	--tag $(REPO):single-$(GIT_SHA) \
		--tag $(REPO):single-$(GIT_BRANCH) \
		--tag $(REPO):single-latest \
		.

.PHONY: image-buildx
image-buildx: Dockerfile $(SOURCE_FILES)
ifneq ($(shell git status --porcelain | wc -l | xargs), 0)
	@$(warning HEAD is not clean, aborting image build)
	@false
endif
	@docker buildx inspect --builder multi || docker buildx create --name multi --use
	@docker buildx build --builder multi \
  --platform linux/amd64,linux/arm/v7,linux/arm64 \
  --cache-to type=inline \
  --label org.opencontainers.image.title=$(OCI_TITLE) \
  --label org.opencontainers.image.description="$(OCI_DESCRIPTION)" \
  --label org.opencontainers.image.url=$(OCI_URL) \
  --label org.opencontainers.image.source=$(OCI_SOURCE) \
  --label org.opencontainers.image.version=$(OCI_VERSION) \
  --label org.opencontainers.image.created=$(OCI_CREATED) \
  --label org.opencontainers.image.revision=$(OCI_REVISION) \
  --label org.opencontainers.image.licenses=$(OCI_LICENSES) \
	--label org.opencontainers.image.authors="$(OCI_AUTHORS)" \
	--label org.opencontainers.image.documentation=$(OCI_DOCUMENTATION) \
	--label org.opencontainers.image.vendor="$(OCI_VENDOR)" \
  --cache-from $(REPO):$(GIT_SHA) \
  --cache-from $(REPO):$(GIT_BRANCH) \
	--cache-from $(REPO):master \
	--cache-from $(REPO):latest \
  --tag $(REPO):$(GIT_SHA) \
  --tag $(REPO):$(GIT_BRANCH) \
  --tag $(REPO):latest \
  --file ./Dockerfile .

.PHONY: image-sh
image-sh: image
	docker run --rm -it --entrypoint=ash wwmoraes/kubegraph:single-latest

.PHONY: release
release:
	env -u GITLAB_TOKEN goreleaser release --rm-dist

.PHONY: test-release
test-release:
	env -u GITLAB_TOKEN goreleaser release --rm-dist --snapshot
