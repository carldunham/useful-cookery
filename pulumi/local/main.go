package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Load configuration
		conf := config.New(ctx, "")
		namespace := conf.Get("namespace")
		if namespace == "" {
			namespace = "useful-cookery-local"
		}

		// Create a provider for the local k3d cluster
		k8sProvider, err := kubernetes.NewProvider(ctx, "k3d-provider", &kubernetes.ProviderArgs{
			Context: pulumi.String("k3d-useful-cookery"),
		})
		if err != nil {
			return err
		}

		// Create namespace
		ns, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(namespace),
				Labels: pulumi.StringMap{
					"name": pulumi.String(namespace),
				},
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		// Create ConfigMap for application configuration
		_, err = corev1.NewConfigMap(ctx, "useful-cookery-config", &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("useful-cookery-config"),
				Namespace: ns.Metadata.Name(),
			},
			Data: pulumi.StringMap{
				"API_URL":      pulumi.String("http://useful-cookery.local/api"),
				"DATABASE_URL": pulumi.String("postgres://postgres:postgres@postgres:5432/useful-cookery?sslmode=disable"),
				"ENV":          pulumi.String("local"),
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		// Create API service
		apiService, err := corev1.NewService(ctx, "useful-cookery-api", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("useful-cookery-api"),
				Namespace: ns.Metadata.Name(),
				Labels: pulumi.StringMap{
					"app":         pulumi.String("useful-cookery-api"),
					"environment": pulumi.String("local"),
				},
			},
			Spec: &corev1.ServiceSpecArgs{
				Selector: pulumi.StringMap{
					"app": pulumi.String("useful-cookery-api"),
				},
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(8080),
						TargetPort: pulumi.Int(8080),
					},
				},
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		// Create UI service
		uiService, err := corev1.NewService(ctx, "useful-cookery-ui", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("useful-cookery-ui"),
				Namespace: ns.Metadata.Name(),
				Labels: pulumi.StringMap{
					"app":         pulumi.String("useful-cookery-ui"),
					"environment": pulumi.String("local"),
				},
			},
			Spec: &corev1.ServiceSpecArgs{
				Selector: pulumi.StringMap{
					"app": pulumi.String("useful-cookery-ui"),
				},
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(80),
					},
				},
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("namespace", ns.Metadata.Name())
		ctx.Export("apiService", apiService.Metadata.Name())
		ctx.Export("uiService", uiService.Metadata.Name())
		ctx.Export("ingressHost", pulumi.String("useful-cookery.local"))

		return nil
	})
}
