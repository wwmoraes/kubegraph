<p align="center">
 <img width=400px height=279px src="https://raw.githubusercontent.com/wwmoraes/kubegraph/master/sample.png" alt="kubegraph sample"></a>
</p>

<h3 align="center">Kubegraph</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)][![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwwmoraes%2Fkubegraph.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwwmoraes%2Fkubegraph?ref=badge_shield)
()
[![GitHub Issues](https://img.shields.io/github/issues/wwmoraes/kubegraph.svg)](https://github.com/wwmoraes/kubegraph/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/wwmoraes/kubegraph.svg)](https://github.com/wwmoraes/kubegraph/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/wwmoraes/kubegraph)](https://hub.docker.com/r/wwmoraes/kubegraph)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/wwmoraes/kubegraph?label=image%20version)](https://hub.docker.com/r/wwmoraes/kubegraph)
[![Docker Pulls](https://img.shields.io/docker/pulls/wwmoraes/kubegraph)](https://hub.docker.com/r/wwmoraes/kubegraph)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=alert_status)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=bugs)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=security_rating)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)

[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=coverage)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=code_smells)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_kubegraph&metric=sqale_index)](https://sonarcloud.io/dashboard?id=wwmoraes_kubegraph)

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

Kubegraph is a CLI tool that parses kubernetes resources and generates a graph
with the relations between those. The graph is done using Graphviz, and can be
further customized after generation.

## 🏁 Getting Started <a name = "getting_started"></a>

Fetch the dependencies and build with

```shell
make build
```

### Prerequisites

Kubegraph is done using Golang 1.15, using a pure Go graphviz implementation to
generate the graph.

Everything is set as a direct dependency, and `go mod vendor` will install for you.

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

or using the docker image

```shell
docker run --rm -it -v ${PWD}:/home/kubegraph wwmoraes/kubegraph:latest my-multidoc.yaml
```

On this example, a `my-multidoc.dot` file will be generated, alongside with an
`icons` folder used by the graph. You can convert it to an image using `dot`, e.g.

```shell
dot -Tsvg -o my-multidoc.svg my-multidoc.dot
```

If your graphviz installation has been compiled with pango, cairo and rsvg, you'll
also be able to generate static formats as png or jpeg. Do note that currently
kubegraph uses svg icons, and cairo mess up when generating raster images with
those (namely they'll either look blurred or won't be drawn at all). A future
version will address this by using raster icons.

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
- [wwmoraes/dot](https://github.com/wwmoraes/dot) - plain Go Graphviz package
- [k8s.io/client-go](https://github.com/kubernetes/client-go) - Kubernetes Go client
- [kubernetes/community](https://github.com/kubernetes/community) - amazing icons
- [spf13/cobra](github.com/spf13/cobra) - CLI framework

## ✍️ Authors <a name = "authors"></a>

- [@wwmoraes](https://github.com/wwmoraes) - Idea & Initial work

## 🎉 Acknowledgements <a name = "acknowledgement"></a>

- Kubernetes sigs members for the excellent abstractions and interfaces available on Golang
- [@damianopetrungaro](https://github.com/damianopetrungaro) for the honest reviews and patience


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwwmoraes%2Fkubegraph.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwwmoraes%2Fkubegraph?ref=badge_large)