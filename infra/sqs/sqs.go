package sqs


import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//Create a sqs client and return
func Connect(sess *session.Session) *sqs.SQS  {
	 var conn *sqs.SQS
	if sess == nil{
		//TODO revisit to check is this load config file details too.
		defaultSvc := session.New()
		conn = sqs.New(defaultSvc)
	} else {
		conn = sqs.New(sess)
	}

	 return conn
}

