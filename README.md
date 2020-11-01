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

## ⛏️ Built Using <a name = "built_using"></a>

- [Golang](https://golang.org) - Base language
- [goccy/go-graphviz](https://github.com/goccy/go-graphviz) - Graphviz C bindings
- [k8s.io/client-go](https://github.com/kubernetes/client-go) - Kubernetes Go client

## ✍️ Authors <a name = "authors"></a>

- [@wwmoraes](https://github.com/wwmoraes) - Idea & Initial work

## 🎉 Acknowledgements <a name = "acknowledgement"></a>

- Kubernetes sigs members for the excellent abstractions and interfaces available on Golang
