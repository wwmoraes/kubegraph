package adapters

import (
	_ "github.com/wwmoraes/kubegraph/internal/adapters/apiextensions"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/app.k8s.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/apps"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/autoscaling"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/batch"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/core"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/extensions"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/networking.istio.io"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/policy"
	_ "github.com/wwmoraes/kubegraph/internal/adapters/rbac"
	// _ "github.com/wwmoraes/kubegraph/internal/adapters/dummy"
)
