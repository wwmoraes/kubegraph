ICONS_FOLDER := icons
ICONS_PKG := icons
ICONS_GO_FILE := $(ICONS_FOLDER)/icons.go
ICONS_FILES := $(ICONS_FOLDER)/*.svg
SOURCE_FILES := $(wildcard cmd/kubegraph/*.go) $(wildcard internal/*/*.go) $(wildcard icons/*)

CGO_ENABLED := 1
CGO_LDFLAGS := -g -O2 -v

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

.PHONY: coverage-html
coverage-html: coverage.html

coverage.html: coverage.out
	go tool cover -html=$< -o $@

coverage.out: $(SOURCE_FILES)
	go test -race -cover -coverprofile=coverage.out -v ./...

build: $(SOURCE_FILES) vendor
	go build -mod=vendor ./...

vendor: go.mod go.sum
	go mod vendor

.PHONY: run
run:
	go run cmd/kubegraph/main.go sample.yaml

.PHONY: image
image:
	docker build -t wwmoraes/kubegraph:latest .