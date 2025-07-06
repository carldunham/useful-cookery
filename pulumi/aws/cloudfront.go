package aws

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudFrontArgs struct {
	Name              string
	DomainName        string
	ALBDomainName     pulumi.StringInput
	ACMCertificateARN string
	Environment       string
}

type CloudFrontOutputs struct {
	Distribution *cloudfront.Distribution
	DomainName   pulumi.StringOutput
}

func NewCloudFront(ctx *pulumi.Context, args *CloudFrontArgs) (*CloudFrontOutputs, error) {
	// We don't need an OAI for ALB origin, but we would for S3 origins

	// Create cache policy
	cachePolicy, err := cloudfront.NewCachePolicy(ctx, args.Name+"-cache-policy", &cloudfront.CachePolicyArgs{
		Comment:    pulumi.String("Cache policy for " + args.Name),
		DefaultTtl: pulumi.Int(86400),
		MaxTtl:     pulumi.Int(31536000),
		MinTtl:     pulumi.Int(1),
		Name:       pulumi.String(args.Name + "-cache-policy"),
		ParametersInCacheKeyAndForwardedToOrigin: &cloudfront.CachePolicyParametersInCacheKeyAndForwardedToOriginArgs{
			CookiesConfig: &cloudfront.CachePolicyParametersInCacheKeyAndForwardedToOriginCookiesConfigArgs{
				CookieBehavior: pulumi.String("none"),
			},
			HeadersConfig: &cloudfront.CachePolicyParametersInCacheKeyAndForwardedToOriginHeadersConfigArgs{
				HeaderBehavior: pulumi.String("none"),
			},
			QueryStringsConfig: &cloudfront.CachePolicyParametersInCacheKeyAndForwardedToOriginQueryStringsConfigArgs{
				QueryStringBehavior: pulumi.String("none"),
			},
			EnableAcceptEncodingGzip:   pulumi.Bool(true),
			EnableAcceptEncodingBrotli: pulumi.Bool(true),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create origin request policy
	originRequestPolicy, err := cloudfront.NewOriginRequestPolicy(ctx, args.Name+"-origin-request-policy", &cloudfront.OriginRequestPolicyArgs{
		Comment: pulumi.String("Origin request policy for " + args.Name),
		CookiesConfig: &cloudfront.OriginRequestPolicyCookiesConfigArgs{
			CookieBehavior: pulumi.String("none"),
		},
		HeadersConfig: &cloudfront.OriginRequestPolicyHeadersConfigArgs{
			HeaderBehavior: pulumi.String("whitelist"),
			Headers: &cloudfront.OriginRequestPolicyHeadersConfigHeadersArgs{
				Items: pulumi.StringArray{
					pulumi.String("Origin"),
					pulumi.String("Host"),
				},
			},
		},
		QueryStringsConfig: &cloudfront.OriginRequestPolicyQueryStringsConfigArgs{
			QueryStringBehavior: pulumi.String("all"),
		},
		Name: pulumi.String(args.Name + "-origin-request-policy"),
	})
	if err != nil {
		return nil, err
	}

	// Create response headers policy
	responseHeadersPolicy, err := cloudfront.NewResponseHeadersPolicy(ctx, args.Name+"-response-headers-policy", &cloudfront.ResponseHeadersPolicyArgs{
		Comment: pulumi.String("Response headers policy for " + args.Name),
		CorsConfig: &cloudfront.ResponseHeadersPolicyCorsConfigArgs{
			AccessControlAllowCredentials: pulumi.Bool(false),
			AccessControlAllowHeaders: &cloudfront.ResponseHeadersPolicyCorsConfigAccessControlAllowHeadersArgs{
				Items: pulumi.StringArray{
					pulumi.String("*"),
				},
			},
			AccessControlAllowMethods: &cloudfront.ResponseHeadersPolicyCorsConfigAccessControlAllowMethodsArgs{
				Items: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
					pulumi.String("OPTIONS"),
					pulumi.String("PUT"),
					pulumi.String("POST"),
					pulumi.String("DELETE"),
					pulumi.String("PATCH"),
				},
			},
			AccessControlAllowOrigins: &cloudfront.ResponseHeadersPolicyCorsConfigAccessControlAllowOriginsArgs{
				Items: pulumi.StringArray{
					pulumi.String("*"),
				},
			},
			AccessControlExposeHeaders: &cloudfront.ResponseHeadersPolicyCorsConfigAccessControlExposeHeadersArgs{
				Items: pulumi.StringArray{
					pulumi.String("*"),
				},
			},
			AccessControlMaxAgeSec: pulumi.Int(600),
			OriginOverride:         pulumi.Bool(true),
		},
		SecurityHeadersConfig: &cloudfront.ResponseHeadersPolicySecurityHeadersConfigArgs{
			ContentSecurityPolicy: &cloudfront.ResponseHeadersPolicySecurityHeadersConfigContentSecurityPolicyArgs{
				ContentSecurityPolicy: pulumi.String("default-src 'self'; img-src 'self' data:; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; connect-src 'self' https://*.amazonaws.com;"),
				Override:              pulumi.Bool(true),
			},
			FrameOptions: &cloudfront.ResponseHeadersPolicySecurityHeadersConfigFrameOptionsArgs{
				FrameOption: pulumi.String("DENY"),
				Override:    pulumi.Bool(true),
			},
			ReferrerPolicy: &cloudfront.ResponseHeadersPolicySecurityHeadersConfigReferrerPolicyArgs{
				ReferrerPolicy: pulumi.String("same-origin"),
				Override:       pulumi.Bool(true),
			},
			StrictTransportSecurity: &cloudfront.ResponseHeadersPolicySecurityHeadersConfigStrictTransportSecurityArgs{
				AccessControlMaxAgeSec: pulumi.Int(63072000),
				IncludeSubdomains:      pulumi.Bool(true),
				Preload:                pulumi.Bool(true),
				Override:               pulumi.Bool(true),
			},
			// XContentTypeOptions field is not available in the current version
			XssProtection: &cloudfront.ResponseHeadersPolicySecurityHeadersConfigXssProtectionArgs{
				Protection: pulumi.Bool(true),
				ModeBlock:  pulumi.Bool(true),
				Override:   pulumi.Bool(true),
			},
		},
		Name: pulumi.String(args.Name + "-response-headers-policy"),
	})
	if err != nil {
		return nil, err
	}

	// Create distribution
	distribution, err := cloudfront.NewDistribution(ctx, args.Name, &cloudfront.DistributionArgs{
		Enabled:           pulumi.Bool(true),
		IsIpv6Enabled:     pulumi.Bool(true),
		HttpVersion:       pulumi.String("http2and3"),
		PriceClass:        pulumi.String("PriceClass_100"),
		Comment:           pulumi.String("Distribution for " + args.Name),
		DefaultRootObject: pulumi.String("index.html"),
		Aliases: pulumi.StringArray{
			pulumi.String(args.DomainName),
		},
		Origins: cloudfront.DistributionOriginArray{
			&cloudfront.DistributionOriginArgs{
				DomainName: args.ALBDomainName,
				OriginId:   pulumi.String("alb"),
				CustomOriginConfig: &cloudfront.DistributionOriginCustomOriginConfigArgs{
					HttpPort:             pulumi.Int(80),
					HttpsPort:            pulumi.Int(443),
					OriginProtocolPolicy: pulumi.String("https-only"),
					OriginSslProtocols: pulumi.StringArray{
						pulumi.String("TLSv1.2"),
					},
				},
			},
		},
		DefaultCacheBehavior: &cloudfront.DistributionDefaultCacheBehaviorArgs{
			AllowedMethods: pulumi.StringArray{
				pulumi.String("GET"),
				pulumi.String("HEAD"),
				pulumi.String("OPTIONS"),
				pulumi.String("PUT"),
				pulumi.String("POST"),
				pulumi.String("PATCH"),
				pulumi.String("DELETE"),
			},
			CachedMethods: pulumi.StringArray{
				pulumi.String("GET"),
				pulumi.String("HEAD"),
				pulumi.String("OPTIONS"),
			},
			TargetOriginId:          pulumi.String("alb"),
			ViewerProtocolPolicy:    pulumi.String("redirect-to-https"),
			CachePolicyId:           cachePolicy.ID(),
			OriginRequestPolicyId:   originRequestPolicy.ID(),
			ResponseHeadersPolicyId: responseHeadersPolicy.ID(),
			Compress:                pulumi.Bool(true),
		},
		ViewerCertificate: &cloudfront.DistributionViewerCertificateArgs{
			AcmCertificateArn:      pulumi.String(args.ACMCertificateARN),
			SslSupportMethod:       pulumi.String("sni-only"),
			MinimumProtocolVersion: pulumi.String("TLSv1.2_2021"),
		},
		Restrictions: &cloudfront.DistributionRestrictionsArgs{
			GeoRestriction: &cloudfront.DistributionRestrictionsGeoRestrictionArgs{
				RestrictionType: pulumi.String("none"),
			},
		},
		Tags: pulumi.StringMap{
			"Name":        pulumi.String(args.Name),
			"Environment": pulumi.String(args.Environment),
		},
	})
	if err != nil {
		return nil, err
	}

	return &CloudFrontOutputs{
		Distribution: distribution,
		DomainName:   distribution.DomainName,
	}, nil
}
