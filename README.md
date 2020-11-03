<p align="center">
 <img width=400px height=279px src="https://raw.githubusercontent.com/wwmoraes/kubegraph/master/sample.png" alt="kubegraph sample"></a>
</p>

<h3 align="center">Kubegraph</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![GitHub Issues](https://img.shields.io/github/issues/wwmoraes/kubegraph.svg)](https://github.com/wwmoraes/kubegraph/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/wwmoraes/kubegraph.svg)](https://github.com/wwmoraes/kubegraph/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

---

<p align="center"> Yet another kubernetes resource graph generator
    <br>
</p>

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Built Using](#built_using)
- [TODO](../TODO.md)
- [Contributing](../CONTRIBUTING.md)
- [Authors](#authors)
- [Acknowledgments](#acknowledgement)

## 🧐 About <a name = "about"></a>

KubeGraph is a CLI tool that parses kubernetes resources and generates a graph
with the relations between those. The graph is done using Graphviz, and can be
further customized after generation.

## 🏁 Getting Started <a name = "getting_started"></a>

Fetch the dependencies and build with

```shell
git clone git@github.com:wwmoraes/go-graphviz.git ../go-graphviz
go mod vendor
go build -mod=vendor ./...
```

### Prerequisites

KubeGraph is done using Golang 1.15, and also depends on a modified version of
[goccy/go-graphviz](https://github.com/goccy/go-graphviz), [wwmoraes/go-graphviz](github.com/wwmoraes/go-graphviz),
while it is not merged into the former. This can be seen on the `go.mod` file as

```text
replace github.com/goccy/go-graphviz => ../go-graphviz
```

Thus why you need to explicitly clone the latter repository before vendoring or
building.

Everything else is set as a direct dependency, and `go mod vendor` will install
for you.

### Installing

It can be installed using standard go install

```shell
go install ./...
```

Then, if you have GOPATH on your path, you can call `kubepath` directly anywhere.

## 🔧 Running the tests <a name = "tests"></a>

WIP, there's no tests yet 😞

## 🎈 Usage <a name="usage"></a>

```shell
kubegraph my-multidoc.yaml
```

### How to add support for a single/suite of custom resource definitions

First, import the scheme and add it to client-go's scheme on `internal/loader/getDecoder.go@init`:

```go
import (
  "k8s.io/client-go/kubernetes/scheme"
  // import the target scheme
  myAwesomeScheme "githost.com/owner/repository/pkg/client/clientset/scheme"
)

func init() {
  // add the scheme to client-go's scheme
  _ = myAwesomeScheme.AddToScheme(scheme.Scheme)
}
```

then:

1. vendor it with `go mod vendor` to update `go.mod` and `go.sum`

1. add adapters for the kinds on that scheme at `internal/adapters/<api-group>/<api-version>`. You can
copy from an existing one, or use the `internal/adapters/dummy/v1/dummy.go` as a guide.

1. import your API versions on the group level (check `internal/adapters/dummy/dummy.go`)

1. import the group on the top level on `internal/adapters/adapters.go`

1. [optional, recommended] add a SVG icon for the new kinds on `icons/` and
set it on your adapter's `Create` function, on the call to `statefulGraph.AddStyledNode`

1. regenerate the icons embedded asset module with `make icons`

1. commit and profit :D

## ⛏️ Built Using <a name = "built_using"></a>

- [Golang](https://golang.org) - Base language
- [goccy/go-graphviz](https://github.com/goccy/go-graphviz) - Graphviz C bindings
- [k8s.io/client-go](https://github.com/kubernetes/client-go) - Kubernetes Go client
- [kubernetes/community](https://github.com/kubernetes/community) - amazing icons
- [spf13/cobra](github.com/spf13/cobra) - CLI framework

## ✍️ Authors <a name = "authors"></a>

- [@wwmoraes](https://github.com/wwmoraes) - Idea & Initial work

## 🎉 Acknowledgements <a name = "acknowledgement"></a>

- Kubernetes sigs members for the excellent abstractions and interfaces available on Golang
- [@damianopetrungaro](https://github.com/damianopetrungaro) for the honest reviews and patience
