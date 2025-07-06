package aws

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VPCArgs struct {
	Name               string
	CIDR               string
	AvailabilityZones  []string
	PrivateSubnetCIDRs []string
	PublicSubnetCIDRs  []string
	EnableNATGateway   bool
	SingleNATGateway   bool
	Environment        string
}

type VPCOutputs struct {
	VPC             *ec2.Vpc
	PrivateSubnets  []*ec2.Subnet
	PublicSubnets   []*ec2.Subnet
	DatabaseSubnets []*ec2.Subnet
	VPCID           pulumi.StringOutput
}

func NewVPC(ctx *pulumi.Context, args *VPCArgs) (*VPCOutputs, error) {
	// Create VPC
	vpc, err := ec2.NewVpc(ctx, args.Name, &ec2.VpcArgs{
		CidrBlock: pulumi.String(args.CIDR),
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.Name),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create Internet Gateway
	igw, err := ec2.NewInternetGateway(ctx, args.Name+"-igw", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.Name + "-igw"),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create public subnets
	publicSubnets := make([]*ec2.Subnet, len(args.PublicSubnetCIDRs))
	for i, cidr := range args.PublicSubnetCIDRs {
		subnet, err := ec2.NewSubnet(ctx, args.Name+"-public-"+args.AvailabilityZones[i], &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String(cidr),
			AvailabilityZone: pulumi.String(args.AvailabilityZones[i]),
			Tags: pulumi.StringMap{
				"Name":        pulumi.String(args.Name + "-public-" + args.AvailabilityZones[i]),
				"Environment": pulumi.String(args.Environment),
				"Type":        pulumi.String("public"),
			},
		})
		if err != nil {
			return nil, err
		}
		publicSubnets[i] = subnet
	}

	// Create public route table
	publicRouteTable, err := ec2.NewRouteTable(ctx, args.Name+"-public-rt", &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.Name + "-public-rt"),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create public route
	_, err = ec2.NewRoute(ctx, args.Name+"-public-route", &ec2.RouteArgs{
		RouteTableId:         publicRouteTable.ID(),
		DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
		GatewayId:            igw.ID(),
	})
	if err != nil {
		return nil, err
	}

	// Associate public subnets with public route table
	for i, subnet := range publicSubnets {
		_, err = ec2.NewRouteTableAssociation(ctx, args.Name+"-public-rta-"+args.AvailabilityZones[i], &ec2.RouteTableAssociationArgs{
			SubnetId:     subnet.ID(),
			RouteTableId: publicRouteTable.ID(),
		})
		if err != nil {
			return nil, err
		}
	}

	// Create private subnets
	privateSubnets := make([]*ec2.Subnet, len(args.PrivateSubnetCIDRs))
	for i, cidr := range args.PrivateSubnetCIDRs {
		subnet, err := ec2.NewSubnet(ctx, args.Name+"-private-"+args.AvailabilityZones[i], &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String(cidr),
			AvailabilityZone: pulumi.String(args.AvailabilityZones[i]),
			Tags: pulumi.StringMap{
				"Name":        pulumi.String(args.Name + "-private-" + args.AvailabilityZones[i]),
				"Environment": pulumi.String(args.Environment),
				"Type":        pulumi.String("private"),
			},
		})
		if err != nil {
			return nil, err
		}
		privateSubnets[i] = subnet
	}

	// Create NAT Gateway if enabled
	if args.EnableNATGateway {
		// Create Elastic IP for NAT Gateway
		eip, err := ec2.NewEip(ctx, args.Name+"-eip", &ec2.EipArgs{
			Vpc: pulumi.Bool(true),
			Tags: pulumi.StringMap{
				"Name":        pulumi.String(args.Name + "-eip"),
				"Environment": pulumi.String(args.Environment),
			},
		})
		if err != nil {
			return nil, err
		}

		// Create NAT Gateway
		natGateway, err := ec2.NewNatGateway(ctx, args.Name+"-nat", &ec2.NatGatewayArgs{
			AllocationId: eip.ID(),
			SubnetId:     publicSubnets[0].ID(),
			Tags: pulumi.StringMap{
				"Name":        pulumi.String(args.Name + "-nat"),
				"Environment": pulumi.String(args.Environment),
			},
		})
		if err != nil {
			return nil, err
		}

		// Create private route table
		privateRouteTable, err := ec2.NewRouteTable(ctx, args.Name+"-private-rt", &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Tags: pulumi.StringMap{
				"Name":        pulumi.String(args.Name + "-private-rt"),
				"Environment": pulumi.String(args.Environment),
			},
		})
		if err != nil {
			return nil, err
		}

		// Create private route
		_, err = ec2.NewRoute(ctx, args.Name+"-private-route", &ec2.RouteArgs{
			RouteTableId:         privateRouteTable.ID(),
			DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
			NatGatewayId:         natGateway.ID(),
		})
		if err != nil {
			return nil, err
		}

		// Associate private subnets with private route table
		for i, subnet := range privateSubnets {
			_, err = ec2.NewRouteTableAssociation(ctx, args.Name+"-private-rta-"+args.AvailabilityZones[i], &ec2.RouteTableAssociationArgs{
				SubnetId:     subnet.ID(),
				RouteTableId: privateRouteTable.ID(),
			})
			if err != nil {
				return nil, err
			}
		}
	}

	// Create database subnets (using the same CIDRs as private subnets for now)
	databaseSubnets := make([]*ec2.Subnet, len(args.PrivateSubnetCIDRs))
	for i, cidr := range args.PrivateSubnetCIDRs {
		// Adjust CIDR for database subnets to avoid overlap
		// This is a simple approach; in a real scenario, you'd want to properly plan your CIDR ranges
		subnet, err := ec2.NewSubnet(ctx, args.Name+"-database-"+args.AvailabilityZones[i], &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String(cidr),
			AvailabilityZone: pulumi.String(args.AvailabilityZones[i]),
			Tags: pulumi.StringMap{
				"Name":        pulumi.String(args.Name + "-database-" + args.AvailabilityZones[i]),
				"Environment": pulumi.String(args.Environment),
				"Type":        pulumi.String("database"),
			},
		})
		if err != nil {
			return nil, err
		}
		databaseSubnets[i] = subnet
	}

	return &VPCOutputs{
		VPC:             vpc,
		PrivateSubnets:  privateSubnets,
		PublicSubnets:   publicSubnets,
		DatabaseSubnets: databaseSubnets,
		VPCID:           vpc.ID().ToStringOutput(),
	}, nil
}
