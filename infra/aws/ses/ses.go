/*
Package ses provides the library to communicate to ses service
*/
package ses

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// New creates a new instance of AWS SES
func New(config *aws.Config) (*ses.SES, error) {
	awsSession := session.New(config)
	sesService := ses.New(awsSession)
	return sesService, nil
}
