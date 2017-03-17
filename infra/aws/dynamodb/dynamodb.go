package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Util defines the structure to provide additional utils to dynamodb
type Util struct {
	FormatTableName func(string) *string
}

// DynamoDB defines a dynamodb instance
type DynamoDB struct {
	*dynamodb.DynamoDB
	*Util
}

// New creates a new instance of AWS dynamodb
func New(config *aws.Config, util *Util) *DynamoDB {
	awsSession := session.New(config)
	return &DynamoDB{
		DynamoDB: dynamodb.New(awsSession),
		Util:     util,
	}
}
