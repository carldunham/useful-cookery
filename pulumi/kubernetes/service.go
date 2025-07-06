package kubernetes

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ServiceArgs struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
	Type        string
	Port        int
	TargetPort  int
	Selector    map[string]string
}

type ServiceOutputs struct {
	Service *corev1.Service
}

func NewService(ctx *pulumi.Context, args *ServiceArgs) (*ServiceOutputs, error) {
	// Create a service
	labels := pulumi.StringMap{}
	for k, v := range args.Labels {
		labels[k] = pulumi.String(v)
	}

	annotations := pulumi.StringMap{}
	for k, v := range args.Annotations {
		annotations[k] = pulumi.String(v)
	}

	selector := pulumi.StringMap{}
	for k, v := range args.Selector {
		selector[k] = pulumi.String(v)
	}

	service, err := corev1.NewService(ctx, args.Name, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:        pulumi.String(args.Name),
			Namespace:   pulumi.String(args.Namespace),
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: &corev1.ServiceSpecArgs{
			Type:     pulumi.String(args.Type),
			Selector: selector,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       pulumi.Int(args.Port),
					TargetPort: pulumi.Int(args.TargetPort),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &ServiceOutputs{
		Service: service,
	}, nil
}
