package v1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/api/core/v1 -n coreV1 -t ConfigMap --icon cm
//go:generate adapter -i k8s.io/api/core/v1 -n coreV1 -t PersistentVolume --icon pv
//go:generate adapter -i k8s.io/api/core/v1 -n coreV1 -t PersistentVolumeClaim --icon pvc
//go:generate adapter -i k8s.io/api/core/v1 -n coreV1 -t Pod --icon pod
//go:generate adapter -i k8s.io/api/core/v1 -n coreV1 -t Secret --icon secret
//go:generate adapter -i k8s.io/api/core/v1 -n coreV1 -t Service --icon svc
//go:generate adapter -i k8s.io/api/core/v1 -n coreV1 -t ServiceAccount --icon sa
//go:generate adapter -i k8s.io/api/apps/v1 -n appsV1 -t StatefulSet --icon sts
