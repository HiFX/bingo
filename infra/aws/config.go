package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Config creates a new AWS configuration for initializing new services
func Config(regionName, accessKey, secretKey string) (*aws.Config, error) {
	sess := session.New()
	awsConfig := aws.NewConfig().WithCredentials(credentials.NewChainCredentials([]credentials.Provider{
		&ec2rolecreds.EC2RoleProvider{
			Client: ec2metadata.New(sess),
		},
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
				SessionToken:    "",
			},
		},
	})).WithRegion(regionName)
	return awsConfig, nil
}
