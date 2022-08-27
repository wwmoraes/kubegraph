package v1beta1

//go:generate -command adapter go run github.com/wwmoraes/kubegraph/cmd/adapter gen
//go:generate adapter -i k8s.io/api/admissionregistration/v1beta1 -n admissionregistrationV1beta1 -t MutatingWebhookConfiguration
//go:generate adapter -i k8s.io/api/admissionregistration/v1beta1 -n admissionregistrationV1beta1 -t ValidatingWebhookConfiguration
