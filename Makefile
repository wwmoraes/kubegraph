ICONS_FOLDER := icons
ICONS_PKG := icons
ICONS_GO_FILE := $(ICONS_FOLDER)/icons.go
ICONS_FILES := $(ICONS_FOLDER)/*.svg
CMD_SOURCE_FILES := $(shell find cmd -type f -name '*.go')
INTERNAL_SOURCE_FILES := $(shell find internal -type f -name '*.go')
WIRE_SRC_FILES := $(shell find internal -type f -name 'wire*.go' -not -name '*_gen.go')
WIRE_GEN_FILES := $(patsubst %.go,%_gen.go,$(WIRE_SRC_FILES))
ICONS_SOURCE_FILES := $(wildcard icons/*.go)
SOURCE_FILES := $(CMD_SOURCE_FILES) $(INTERNAL_SOURCE_FILES) $(ICONS_SOURCE_FILES)
GRAPHS_FOLDER := docs
GRAPHS_SRC_FILES := $(wildcard $(GRAPHS_FOLDER)/*.puml)
GRAPHS_SVG_FILES := $(patsubst %.puml,%.svg,$(GRAPHS_SRC_FILES))
GRAPHS_PNG_FILES := $(patsubst %.puml,%.png,$(GRAPHS_SRC_FILES))
YAMLS := $(shell find . \( -name "*.yaml" -o -name "*.yml" \) -not -path "./vendor/*")

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
OCI_VENDOR = $(USERNAME) <$(EMAIL)>
OCI_DOCUMENTATION = https://github.com/$(REPO)

.DEFAULT_GOAL := build

TMPDIR ?= $(or $(TMPDIR),$(shell dirname $(shell mktemp -u)))

.PHONY: icons
icons: $(ICONS_GO_FILE)

$(ICONS_GO_FILE): $(ICONS_FILES)
	@type -p go-bindata 1>/dev/null || go get -u github.com/go-bindata/go-bindata/...
	$(info removing old icons package file...)
	-@rm $(ICONS_GO_FILE)
	$(info regenerating icons package file...)
	@go-bindata -o $(ICONS_GO_FILE) -pkg $(ICONS_PKG) -nometadata -mode 0664 -nomemcopy $(ICONS_FOLDER)

.PHONY: lint
lint: yamllint
	golangci-lint run

yamllint:
	@yamllint $(YAMLS)

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

.PHONY: build
build: kubegraph

.PHONY: wire
wire: $(WIRE_GEN_FILES)

kubegraph: $(SOURCE_FILES) vendor
	go build -mod=vendor -race -o ./ ./...

vendor: go.sum

go.sum: go.mod
	go mod vendor

.PHONY: run
run: $(SOURCE_FILES) vendor
	go run cmd/kubegraph/main.go sample.yaml
	dot -Tsvg -o sample.svg sample.dot

.PHONY: image
image: GIT_BRANCH_SLUG=$(subst /,-,${GIT_BRANCH})
image: Dockerfile $(SOURCE_FILES)
	@time docker buildx build \
		--cache-to type=local,mode=max,dest=$(TMPDIR)/.buildx-cache/$(REPO) \
		--cache-from type=local,src=$(TMPDIR)/.buildx-cache/$(REPO) \
		--cache-from $(REPO):single-$(GIT_SHA) \
		--cache-from $(REPO):single-$(GIT_BRANCH_SLUG) \
		--cache-from $(REPO):single-master \
		--cache-from $(REPO):single-latest \
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
  	--tag $(REPO):single-$(GIT_SHA) \
		--tag $(REPO):single-$(GIT_BRANCH_SLUG) \
		--tag $(REPO):single-latest \
		--load \
		.

.PHONY: image-buildx
image-buildx: GIT_BRANCH_SLUG=$(subst /,-,${GIT_BRANCH})
image-buildx: Dockerfile $(SOURCE_FILES)
	@docker buildx inspect --builder multi || docker buildx create --name multi --use
	@time docker buildx build --builder multi \
  --platform linux/amd64,linux/arm/v7,linux/arm64 \
	--cache-to type=local,mode=max,dest=$(TMPDIR)/.buildx-cache/$(REPO) \
	--cache-from type=local,src=$(TMPDIR)/.buildx-cache/$(REPO) \
  --cache-from $(REPO):$(GIT_SHA) \
  --cache-from $(REPO):$(GIT_BRANCH_SLUG) \
	--cache-from $(REPO):master \
	--cache-from $(REPO):latest \
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
  --tag $(REPO):$(GIT_SHA) \
  --tag $(REPO):$(GIT_BRANCH_SLUG) \
  --tag $(REPO):latest \
  --file ./Dockerfile .

.PHONY: docs
docs: $(GRAPHS_SVG_FILES) $(GRAPHS_PNG_FILES)

.PHONY: graphs
graphs: $(GRAPHS_FOLDER)/full-gen.puml $(GRAPHS_FOLDER)/core-gen.puml
	@$(MAKE) docs

.PHONY: image-sh
image-sh: image
	docker run --rm -it --entrypoint=ash wwmoraes/kubegraph:single-latest

.PHONY: release
release:
	env -u GITLAB_TOKEN goreleaser release --rm-dist

.PHONY: test-release
test-release:
	env -u GITLAB_TOKEN goreleaser release --rm-dist --snapshot

%_gen.go: %.go
	wire ./...

$(GRAPHS_FOLDER)/full-gen.puml: cmd internal icons
	$(info generating $@...)
	@goplantuml -recursive $^ > $@

$(GRAPHS_FOLDER)/core-gen.puml: internal/registry internal/kubegraph internal/utils
	$(info generating $@...)
	@goplantuml -recursive $^ > $@

$(GRAPHS_FOLDER)/%.svg: $(GRAPHS_FOLDER)/%.puml
	$(info generating $@ from $<...)
	@plantuml -tsvg $<

$(GRAPHS_FOLDER)/%.png: $(GRAPHS_FOLDER)/%.puml
	$(info generating $@ from $<...)
	@plantuml -tpng $<
