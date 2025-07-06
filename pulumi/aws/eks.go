package aws

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type EKSArgs struct {
	ClusterName    string
	ClusterVersion string
	VPCID          pulumi.StringInput
	SubnetIDs      pulumi.StringArrayInput
	NodeGroups     map[string]NodeGroupArgs
	Environment    string
}

type NodeGroupArgs struct {
	DesiredCapacity int
	MinCapacity     int
	MaxCapacity     int
	InstanceTypes   []string
	DiskSize        int
}

type EKSOutputs struct {
	Cluster             *eks.Cluster
	NodeGroups          map[string]*eks.NodeGroup
	ClusterName         pulumi.StringOutput
	ClusterEndpoint     pulumi.StringOutput
	ClusterCA           pulumi.StringPtrOutput
	NodeSecurityGroupID pulumi.StringPtrOutput
	KubeConfigRaw       pulumi.StringOutput
	ALBHostname         pulumi.StringOutput
}

func NewEKS(ctx *pulumi.Context, args *EKSArgs) (*EKSOutputs, error) {
	// Create IAM role for EKS cluster
	clusterRole, err := iam.NewRole(ctx, args.ClusterName+"-cluster-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {
					"Service": "eks.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}`),
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.ClusterName + "-cluster-role"),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Attach policies to cluster role
	_, err = iam.NewRolePolicyAttachment(ctx, args.ClusterName+"-cluster-policy", &iam.RolePolicyAttachmentArgs{
		Role:      clusterRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, args.ClusterName+"-service-policy", &iam.RolePolicyAttachmentArgs{
		Role:      clusterRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSServicePolicy"),
	})
	if err != nil {
		return nil, err
	}

	// Create EKS cluster
	cluster, err := eks.NewCluster(ctx, args.ClusterName, &eks.ClusterArgs{
		RoleArn: clusterRole.Arn,
		Version: pulumi.String(args.ClusterVersion),
		VpcConfig: &eks.ClusterVpcConfigArgs{
			SubnetIds:             args.SubnetIDs,
			EndpointPrivateAccess: pulumi.Bool(true),
			EndpointPublicAccess:  pulumi.Bool(true),
		},
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.ClusterName),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create IAM role for node groups
	nodeRole, err := iam.NewRole(ctx, args.ClusterName+"-node-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {
					"Service": "ec2.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}`),
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.ClusterName + "-node-role"),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Attach policies to node role
	_, err = iam.NewRolePolicyAttachment(ctx, args.ClusterName+"-node-worker-policy", &iam.RolePolicyAttachmentArgs{
		Role:      nodeRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, args.ClusterName+"-node-cni-policy", &iam.RolePolicyAttachmentArgs{
		Role:      nodeRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, args.ClusterName+"-node-registry-policy", &iam.RolePolicyAttachmentArgs{
		Role:      nodeRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"),
	})
	if err != nil {
		return nil, err
	}

	// Create node groups
	nodeGroups := make(map[string]*eks.NodeGroup)
	for name, ngArgs := range args.NodeGroups {
		nodeGroup, err := eks.NewNodeGroup(ctx, args.ClusterName+"-"+name, &eks.NodeGroupArgs{
			ClusterName:   cluster.Name,
			NodeRoleArn:   nodeRole.Arn,
			SubnetIds:     args.SubnetIDs,
			InstanceTypes: pulumi.ToStringArray(ngArgs.InstanceTypes),
			ScalingConfig: &eks.NodeGroupScalingConfigArgs{
				DesiredSize: pulumi.Int(ngArgs.DesiredCapacity),
				MinSize:     pulumi.Int(ngArgs.MinCapacity),
				MaxSize:     pulumi.Int(ngArgs.MaxCapacity),
			},
			DiskSize: pulumi.Int(ngArgs.DiskSize),
			Tags: pulumi.StringMap{
				"Name":        pulumi.String(args.ClusterName + "-" + name),
				"Environment": pulumi.String(args.Environment),
			},
		})
		if err != nil {
			return nil, err
		}
		nodeGroups[name] = nodeGroup
	}

	// Generate kubeconfig
	kubeconfig := pulumi.All(cluster.Name, cluster.Endpoint, cluster.CertificateAuthority.Data()).ApplyT(
		func(args []interface{}) string {
			name := args[0].(string)
			endpoint := args[1].(string)
			ca := args[2].(string)
			return `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ` + ca + `
    server: ` + endpoint + `
  name: ` + name + `
contexts:
- context:
    cluster: ` + name + `
    user: ` + name + `
  name: ` + name + `
current-context: ` + name + `
kind: Config
preferences: {}
users:
- name: ` + name + `
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: aws
      args:
      - eks
      - get-token
      - --cluster-name
      - ` + name + `
`
		}).(pulumi.StringOutput)

	// For now, we'll return an empty ALB hostname
	// In a real implementation, this would be populated after setting up the AWS Load Balancer Controller
	albHostname := pulumi.String("").ToStringOutput()

	return &EKSOutputs{
		Cluster:             cluster,
		NodeGroups:          nodeGroups,
		ClusterName:         cluster.Name,
		ClusterEndpoint:     cluster.Endpoint,
		ClusterCA:           cluster.CertificateAuthority.Data(),
		NodeSecurityGroupID: cluster.VpcConfig.ClusterSecurityGroupId(),
		KubeConfigRaw:       kubeconfig,
		ALBHostname:         albHostname,
	}, nil
}
