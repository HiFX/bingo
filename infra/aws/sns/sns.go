/*
Package sns provides the library to communicate to sns service
*/
package sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// New creates a new instance of AWS SNS
func New(config *aws.Config) (*sns.SNS, error) {
	awsSession := session.New(config)
	snsService := sns.New(awsSession)
	return snsService, nil
}
