package aws

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type RDSArgs struct {
	Identifier              string
	Engine                  string
	EngineVersion           string
	InstanceClass           string
	AllocatedStorage        int
	MaxAllocatedStorage     int
	DBName                  string
	Username                string
	Password                string
	VPCID                   pulumi.StringInput
	SubnetIDs               pulumi.StringArrayInput
	MultiAZ                 bool
	BackupRetentionPeriod   int
	DeletionProtection      bool
	AllowedSecurityGroupIDs []pulumi.StringInput
	Environment             string
}

type RDSOutputs struct {
	Instance *rds.Instance
	Endpoint pulumi.StringOutput
	Port     pulumi.IntOutput
}

func NewRDS(ctx *pulumi.Context, args *RDSArgs) (*RDSOutputs, error) {
	// Create DB subnet group
	subnetGroup, err := rds.NewSubnetGroup(ctx, args.Identifier+"-subnet-group", &rds.SubnetGroupArgs{
		Name:      pulumi.String(args.Identifier + "-subnet-group"),
		SubnetIds: args.SubnetIDs,
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.Identifier + "-subnet-group"),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create security group for RDS
	securityGroup, err := ec2.NewSecurityGroup(ctx, args.Identifier+"-sg", &ec2.SecurityGroupArgs{
		VpcId:       args.VPCID,
		Description: pulumi.String("Security group for RDS instance " + args.Identifier),
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.Identifier + "-sg"),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Allow PostgreSQL traffic from allowed security groups
	for i, sgID := range args.AllowedSecurityGroupIDs {
		_, err = ec2.NewSecurityGroupRule(ctx, fmt.Sprintf("%s-sg-rule-%d", args.Identifier, i), &ec2.SecurityGroupRuleArgs{
			Type:                  pulumi.String("ingress"),
			FromPort:              pulumi.Int(5432),
			ToPort:                pulumi.Int(5432),
			Protocol:              pulumi.String("tcp"),
			SecurityGroupId:       securityGroup.ID(),
			SourceSecurityGroupId: sgID,
		})
		if err != nil {
			return nil, err
		}
	}

	// Create parameter group
	parameterGroup, err := rds.NewParameterGroup(ctx, args.Identifier+"-pg", &rds.ParameterGroupArgs{
		Family: pulumi.String("postgres" + args.EngineVersion),
		Parameters: rds.ParameterGroupParameterArray{
			&rds.ParameterGroupParameterArgs{
				Name:  pulumi.String("max_connections"),
				Value: pulumi.String("100"),
			},
			&rds.ParameterGroupParameterArgs{
				Name:  pulumi.String("shared_buffers"),
				Value: pulumi.String("{DBInstanceClassMemory/32768}"),
			},
		},
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.Identifier + "-pg"),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create RDS instance
	instance, err := rds.NewInstance(ctx, args.Identifier, &rds.InstanceArgs{
		Identifier:            pulumi.String(args.Identifier),
		Engine:                pulumi.String(args.Engine),
		EngineVersion:         pulumi.String(args.EngineVersion),
		InstanceClass:         pulumi.String(args.InstanceClass),
		AllocatedStorage:      pulumi.Int(args.AllocatedStorage),
		MaxAllocatedStorage:   pulumi.Int(args.MaxAllocatedStorage),
		DbName:                pulumi.String(args.DBName),
		Username:              pulumi.String(args.Username),
		Password:              pulumi.String(args.Password),
		DbSubnetGroupName:     subnetGroup.Name,
		VpcSecurityGroupIds:   pulumi.StringArray{securityGroup.ID()},
		ParameterGroupName:    parameterGroup.Name,
		MultiAz:               pulumi.Bool(args.MultiAZ),
		BackupRetentionPeriod: pulumi.Int(args.BackupRetentionPeriod),
		DeletionProtection:    pulumi.Bool(args.DeletionProtection),
		SkipFinalSnapshot:     pulumi.Bool(true),
		StorageEncrypted:      pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.Identifier),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	return &RDSOutputs{
		Instance: instance,
		Endpoint: instance.Endpoint,
		Port:     instance.Port,
	}, nil
}
