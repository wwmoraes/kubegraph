package adapters

import (
	_ "github.com/wwmoraes/kubegraph/internal/adapters/admissionregistration.k8s.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/apiextensions"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/apiregistration.k8s.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/app.k8s.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/apps"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/autoscaling"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/batch"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/core"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/extensions"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/networking.istio.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/policy"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/rbac"
)

//go:generate -command scheme go run github.com/wwmoraes/kubegraph/cmd/adapter scheme
//go:generate scheme -n apiExtensions -i k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme
//go:generate scheme -n application -i sigs.k8s.io/application/api/v1beta1
//go:generate scheme -n aggregator -i k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/scheme
//go:generate scheme -n istio -i istio.io/client-go/pkg/clientset/versioned/scheme

// add any extra schemes here
