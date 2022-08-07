package v1beta1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1 -n apiExtensionsV1beta1 -t CustomResourceDefinition --icon crd
