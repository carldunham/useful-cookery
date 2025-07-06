package kubernetes

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DeploymentArgs struct {
	Name          string
	Namespace     string
	Labels        map[string]string
	Annotations   map[string]string
	Replicas      int
	Image         string
	Port          int
	CPURequest    string
	CPULimit      string
	MemoryRequest string
	MemoryLimit   string
}

type DeploymentOutputs struct {
	Deployment *appsv1.Deployment
}

func NewDeployment(ctx *pulumi.Context, args *DeploymentArgs) (*DeploymentOutputs, error) {
	// Create a simplified deployment with a single container
	labels := pulumi.StringMap{}
	for k, v := range args.Labels {
		labels[k] = pulumi.String(v)
	}

	annotations := pulumi.StringMap{}
	for k, v := range args.Annotations {
		annotations[k] = pulumi.String(v)
	}

	// Create deployment
	deployment, err := appsv1.NewDeployment(ctx, args.Name, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:        pulumi.String(args.Name),
			Namespace:   pulumi.String(args.Namespace),
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(args.Replicas),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: labels,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String(args.Name),
							Image: pulumi.String(args.Image),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(args.Port),
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"cpu":    pulumi.String(args.CPURequest),
									"memory": pulumi.String(args.MemoryRequest),
								},
								Limits: pulumi.StringMap{
									"cpu":    pulumi.String(args.CPULimit),
									"memory": pulumi.String(args.MemoryLimit),
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &DeploymentOutputs{
		Deployment: deployment,
	}, nil
}
