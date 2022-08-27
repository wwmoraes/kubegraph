package v2beta1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/api/autoscaling/v2beta1 -n autoscalingV2beta1 -t HorizontalPodAutoscaler --icon hpa

//go:generate -command dependency go run github.com/wwmoraes/kubegraph/cmd/adapter dep
//go:generate dependency -i k8s.io/api/apps/v1 -n appsV1 -t Deployment
