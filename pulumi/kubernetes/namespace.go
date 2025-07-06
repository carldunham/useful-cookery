package kubernetes

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type NamespaceArgs struct {
	Name        string
	Labels      map[string]string
	Annotations map[string]string
}

type NamespaceOutputs struct {
	Namespace *corev1.Namespace
}

func NewNamespace(ctx *pulumi.Context, args *NamespaceArgs) (*NamespaceOutputs, error) {
	// Create namespace
	namespace, err := corev1.NewNamespace(ctx, args.Name, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:        pulumi.String(args.Name),
			Labels:      pulumi.ToStringMap(args.Labels),
			Annotations: pulumi.ToStringMap(args.Annotations),
		},
	})
	if err != nil {
		return nil, err
	}

	return &NamespaceOutputs{
		Namespace: namespace,
	}, nil
}
