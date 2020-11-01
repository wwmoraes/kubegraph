ICONS_FOLDER := icons
ICONS_PKG := icons
ICONS_GO_FILE := $(ICONS_FOLDER)/icons.go
ICONS_FILES := $(ICONS_FOLDER)/*.svg
SOURCE_FILES := $(wildcard cmd/kubegraph/*.go) $(wildcard internal/*/*.go)

.PHONY: icons
icons: $(ICONS_GO_FILE)

$(ICONS_GO_FILE): $(ICONS_FILES)
	@type -p go-bindata 1>/dev/null || go get -u github.com/go-bindata/go-bindata/...
	$(info removing old icons package file...)
	-@rm $(ICONS_GO_FILE)
	$(info regenerating icons package file...)
	@go-bindata -o $(ICONS_GO_FILE) -pkg $(ICONS_PKG) -nometadata -nomemcopy $(ICONS_FOLDER)

build: $(SOURCE_FILES) vendor
	go build -mod=vendor ./...
	go build ./...

vendor: go.mod go.sum ../go-graphviz
	go mod vendor

../go-graphviz:
	git clone git@github.com:wwmoraes/go-graphviz.git ../go-graphviz

.PHONY: run
run:
	go run cmd/kubegraph/main.go sample.yaml