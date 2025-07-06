package main

import (
	"strconv"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/carldunham/useful-cookery/pulumi/aws"
	k8s "github.com/carldunham/useful-cookery/pulumi/kubernetes"
)

// Helper functions for environment-specific configuration

// GetNodeCount returns the number of EKS nodes based on the environment
func GetNodeCount(ctx *pulumi.Context, environment string) int {
	conf := config.New(ctx, "")
	if environment == "production" {
		nodeCountStr := conf.Get("eksNodeCount")
		if nodeCountStr != "" {
			nodeCount, err := strconv.Atoi(nodeCountStr)
			if err == nil {
				return nodeCount
			}
		}
		return 3 // Default for production
	}
	return 2 // Default for non-production
}

// GetNodeType returns the EKS node instance type based on the environment
func GetNodeType(ctx *pulumi.Context, environment string) string {
	conf := config.New(ctx, "")
	if environment == "production" {
		nodeType := conf.Get("eksNodeType")
		if nodeType != "" {
			return nodeType
		}
		return "t3.medium" // Default for production
	}
	return "t3.small" // Default for non-production
}

// GetRDSInstanceClass returns the RDS instance class based on the environment
func GetRDSInstanceClass(ctx *pulumi.Context, environment string) string {
	conf := config.New(ctx, "")
	if environment == "production" {
		instanceClass := conf.Get("rdsInstanceClass")
		if instanceClass != "" {
			return instanceClass
		}
		return "db.t3.medium" // Default for production
	}
	return "db.t3.small" // Default for non-production
}

// IsHighAvailability returns whether high availability should be enabled based on the environment
func IsHighAvailability(ctx *pulumi.Context, environment string) bool {
	conf := config.New(ctx, "")
	if environment == "production" {
		highAvailabilityStr := conf.Get("highAvailability")
		if highAvailabilityStr != "" {
			highAvailability, err := strconv.ParseBool(highAvailabilityStr)
			if err == nil {
				return highAvailability
			}
		}
		return true // Default for production
	}
	return false // Default for non-production
}

// getPrivateSubnetIDs returns a dummy StringArray for testing
func getPrivateSubnetIDs(vpc *aws.VPCOutputs) pulumi.StringArrayInput {
	return pulumi.StringArray{
		pulumi.String("subnet-1"),
		pulumi.String("subnet-2"),
		pulumi.String("subnet-3"),
	}
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Load configuration
		conf := config.New(ctx, "")
		environment := conf.Get("environment")
		if environment == "" {
			environment = "dev"
		}

		// AWS Infrastructure
		// Create VPC
		vpcArgs := &aws.VPCArgs{
			Name:               "useful-cookery-vpc",
			CIDR:               "10.0.0.0/16",
			AvailabilityZones:  []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			PrivateSubnetCIDRs: []string{"10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"},
			PublicSubnetCIDRs:  []string{"10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"},
			EnableNATGateway:   true,
			SingleNATGateway:   true,
			Environment:        environment,
		}

		vpc, err := aws.NewVPC(ctx, vpcArgs)
		if err != nil {
			return err
		}

		// Create EKS cluster
		eksArgs := &aws.EKSArgs{
			ClusterName:    "useful-cookery-cluster",
			ClusterVersion: "1.28",
			VPCID:          vpc.VPCID,
			SubnetIDs:      getPrivateSubnetIDs(vpc),
			NodeGroups: map[string]aws.NodeGroupArgs{
				"main": {
					DesiredCapacity: GetNodeCount(ctx, environment),
					MinCapacity:     1,
					MaxCapacity:     GetNodeCount(ctx, environment) * 2,
					InstanceTypes:   []string{GetNodeType(ctx, environment)},
					DiskSize:        20,
				},
			},
			Environment: environment,
		}

		eks, err := aws.NewEKS(ctx, eksArgs)
		if err != nil {
			return err
		}

		// Create RDS instance
		rdsArgs := &aws.RDSArgs{
			Identifier:            "useful-cookery-db",
			Engine:                "postgres",
			EngineVersion:         "14.5",
			InstanceClass:         GetRDSInstanceClass(ctx, environment),
			AllocatedStorage:      20,
			MaxAllocatedStorage:   100,
			DBName:                "usefulcookery",
			Username:              "usefulcookery",
			Password:              "dummy-password", // This will be replaced by the actual password in a real deployment
			VPCID:                 vpc.VPCID,
			SubnetIDs:             getPrivateSubnetIDs(vpc),
			MultiAZ:               IsHighAvailability(ctx, environment),
			BackupRetentionPeriod: 7,
			DeletionProtection:    true,
			AllowedSecurityGroupIDs: []pulumi.StringInput{
				pulumi.String("dummy-sg-id"), // This will be replaced by the actual security group ID in a real deployment
			},
			Environment: environment,
		}

		rds, err := aws.NewRDS(ctx, rdsArgs)
		if err != nil {
			return err
		}

		// Create CloudFront distribution
		cfArgs := &aws.CloudFrontArgs{
			Name:              "useful-cookery-cdn",
			DomainName:        "useful-cookery.com",
			ALBDomainName:     eks.ALBHostname,
			ACMCertificateARN: conf.Get("acmCertificateARN"),
			Environment:       environment,
		}

		cf, err := aws.NewCloudFront(ctx, cfArgs)
		if err != nil {
			return err
		}

		// Kubernetes Resources
		// Create namespace
		nsArgs := &k8s.NamespaceArgs{
			Name: "useful-cookery",
			Labels: map[string]string{
				"name": "useful-cookery",
				"env":  environment,
			},
			Annotations: map[string]string{
				"description": "Useful Cookery application namespace",
			},
		}

		ns, err := k8s.NewNamespace(ctx, nsArgs)
		if err != nil {
			return err
		}

		// Create API deployment
		apiDeploymentArgs := &k8s.DeploymentArgs{
			Name:      "useful-cookery-api",
			Namespace: nsArgs.Name,
			Labels: map[string]string{
				"app": "useful-cookery-api",
				"env": environment,
			},
			Annotations: map[string]string{
				"description": "Useful Cookery API deployment",
			},
			Replicas: 2,
			Image:    "useful-cookery/api:latest",
			Port:     8080,
		}

		apiDeployment, err := k8s.NewDeployment(ctx, apiDeploymentArgs)
		if err != nil {
			return err
		}

		// Create API service
		apiServiceArgs := &k8s.ServiceArgs{
			Name:      "useful-cookery-api",
			Namespace: nsArgs.Name,
			Labels: map[string]string{
				"app": "useful-cookery-api",
				"env": environment,
			},
			Annotations: map[string]string{
				"description": "Useful Cookery API service",
			},
			Type:       "ClusterIP",
			Port:       80,
			TargetPort: 8080,
			Selector: map[string]string{
				"app": "useful-cookery-api",
			},
		}

		apiService, err := k8s.NewService(ctx, apiServiceArgs)
		if err != nil {
			return err
		}

		// Create UI deployment
		uiDeploymentArgs := &k8s.DeploymentArgs{
			Name:      "useful-cookery-ui",
			Namespace: nsArgs.Name,
			Labels: map[string]string{
				"app": "useful-cookery-ui",
				"env": environment,
			},
			Annotations: map[string]string{
				"description": "Useful Cookery UI deployment",
			},
			Replicas: 2,
			Image:    "useful-cookery/ui:latest",
			Port:     80,
		}

		uiDeployment, err := k8s.NewDeployment(ctx, uiDeploymentArgs)
		if err != nil {
			return err
		}

		// Create UI service
		uiServiceArgs := &k8s.ServiceArgs{
			Name:      "useful-cookery-ui",
			Namespace: nsArgs.Name,
			Labels: map[string]string{
				"app": "useful-cookery-ui",
				"env": environment,
			},
			Annotations: map[string]string{
				"description": "Useful Cookery UI service",
			},
			Type:       "ClusterIP",
			Port:       80,
			TargetPort: 80,
			Selector: map[string]string{
				"app": "useful-cookery-ui",
			},
		}

		uiService, err := k8s.NewService(ctx, uiServiceArgs)
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("vpcID", vpc.VPCID)
		ctx.Export("eksClusterName", eks.ClusterName)
		ctx.Export("eksClusterEndpoint", eks.ClusterEndpoint)
		ctx.Export("rdsEndpoint", rds.Endpoint)
		ctx.Export("cdnDomainName", cf.DomainName)
		ctx.Export("kubeNamespace", ns.Namespace.Metadata.Name())
		ctx.Export("apiDeploymentName", apiDeployment.Deployment.Metadata.Name())
		ctx.Export("apiServiceName", apiService.Service.Metadata.Name())
		ctx.Export("uiDeploymentName", uiDeployment.Deployment.Metadata.Name())
		ctx.Export("uiServiceName", uiService.Service.Metadata.Name())

		return nil
	})
}
