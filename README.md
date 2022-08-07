# Kubegraph

> Yet another kubernetes resource graph generator

![Kubegraph sample](https://raw.githubusercontent.com/wwmoraes/kubegraph/master/sample.png)

![Status](https://img.shields.io/badge/status-active-success.svg)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/wwmoraes/kubegraph/master.svg)](https://results.pre-commit.ci/latest/github/wwmoraes/kubegraph/master)
[![GitHub Issues](https://img.shields.io/github/issues/wwmoraes/kubegraph.svg)](https://github.com/wwmoraes/kubegraph/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/wwmoraes/kubegraph.svg)](https://github.com/wwmoraes/kubegraph/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwwmoraes%2Fkubegraph.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwwmoraes%2Fkubegraph?ref=badge_shield)

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

---

## 📝 Table of Contents

- [About](#-about)
- [Getting Started](#-getting-started)
- [Usage](#-usage)
- [Built Using](#-built-using)
- [TODO](./TODO.md)
- [Contributing](./CONTRIBUTING.md)
- [Authors](#-authors)
- [Acknowledgments](#-acknowledgements)

## 🧐 About

Kubegraph is a CLI tool that parses kubernetes resources and generates a graph
with the relations between those. The graph is done using Graphviz, and can be
further customized after generation.

## 🏁 Getting Started

Fetch the dependencies and build with

```shell
make build
```

### Prerequisites

- Golang 1.15+
- Graphviz (optional, not used directly by kubegraph)

Kubegraph is written using Golang, and depends on a pure Go graphviz
implementation to generate the graph.

Graphviz can then be used to convert the generated graph to more
portable/conventional formats, such as PNG or SVG.

### Installing

It can be installed using standard go install

```shell
go install github.com/wwmoraes/kubegraph/cmd/kubegraph
# or if you cloned the repository locally
go install ./cmd/kubegraph/...
```

You can then use `kubegraph` if you have GOPATH on your system path.

## 🔧 Running the tests

WIP, there's no tests yet 😞

## 🎈 Usage

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
kubegraph uses svg icons, and cairo messes up when generating raster images with
those (namely they'll either look blurred or won't be drawn at all). A future
version may address this by using raster icons.

### How to add support for a single/suite of custom resource definitions

If adding support for a non-core resource type, then declare the scheme on
`internal/adapters/adapters.go`:

```go
import (
  ...

  //go:generate scheme -n importName -i url/to/scheme
)
```

Run `go generate` to create the scheme import and composition injection code.

Afterwards:

1. vendor it with `go mod vendor` to update `go.mod` and `go.sum`

1. create the folder structure for the target API group and API version:
`internal/adapters/<api-group>/<api-version>`

1. add an anonymous import on the group level, e.g.:

`import _ "github.com/wwmoraes/kubegraph/internal/adapters/<api-group>/<api-version>"`

1. add the adapter generator comments for all resource types you want to add
support for, e.g.:

```go
//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i <import-URL> -n <import-name> -t <kind>
```

1. re-run `go generate` to create the adapters

1. if the new kind(s) have relations with other kinds, then add a `<kind>.go`
file for each, and overload the `Adapter.Configure` method, e.g.:

```go
func (this *<kind>Adapter) Configure(graph registry.StatefulGraph) error {
  <dependency>Adapter, err := Get<dependency>Adapter()
  if err != nil { ... }

  objects, err := this.GetGraphObjects(graph)
  if err != nil { ... }

  // resource is a pointer to the target kubernetes kind struct
  for name, resource := range objects {
    node, err := this.GetGraphNode(graph, name)
    if err != nil { ... }

    // implement here the logic to check if this resource should connect to the
    // dependency, or else skip it

    // connect to the dependency
    targetName := resource.PathTo.Dependency.Name
    _, err := <dependency>Adapter.Connect(graph, node, targetName)
    if err != nil { ... }
  }
}
```

1. [optional, recommended] add a SVG icon for the new kinds on `icons/`

    1. regenerate the icons embedded asset module with `make icons`

    1. add a `--icon <name>` to the adapter generate comment to use it

    1. regenerate the adapters with `go generate`

1. commit and profit :D

## 🛠 Built Using

- [Golang](https://golang.org) - Base language
- [wwmoraes/dot](https://github.com/wwmoraes/dot) - plain Go Graphviz package
- [k8s.io/client-go](https://github.com/kubernetes/client-go) - Kubernetes Go client
- [kubernetes/community](https://github.com/kubernetes/community) - amazing icons
- [spf13/cobra](github.com/spf13/cobra) - CLI framework

## 🧑‍💻 Authors

- [@wwmoraes](https://github.com/wwmoraes) - Idea & Initial work

## 🎉 Acknowledgements

- Kubernetes sigs members for the excellent abstractions and interfaces
available on Golang
- [@damianopetrungaro](https://github.com/damianopetrungaro) for the honest
reviews and patience
