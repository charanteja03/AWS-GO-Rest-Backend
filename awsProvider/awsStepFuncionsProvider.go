package awsprovider

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

//go:generate mockery --output=../mocks --name=AwsStepFunctionInterface|AwsStepFunctionsProvider

//AwsStepFunctionInterface - interface for aws step functions functionality
type AwsStepFunctionInterface interface {
	ListExecutions(input *sfn.ListExecutionsInput) (*sfn.ListExecutionsOutput, error)
	ListStateMachines(input *sfn.ListStateMachinesInput) (*sfn.ListStateMachinesOutput, error)
	DescribeExecution(input *sfn.DescribeExecutionInput) (*sfn.DescribeExecutionOutput, error)
	StartExecution(input *sfn.StartExecutionInput) (*sfn.StartExecutionOutput, error)
}

//AwsStepFunctionsProvider - provider for step function interface
type AwsStepFunctionsProvider interface {
	New(sess *session.Session) (AwsStepFunctionInterface, error)
}

//AwsStepFunctionsRealProvider - struct represents real AWS provider
type AwsStepFunctionsRealProvider struct {
}

//New - function creates a new instance of the SFN client with a session.
func (awsProvider *AwsStepFunctionsRealProvider) New(sess *session.Session) (AwsStepFunctionInterface, error) {
	return sfn.New(sess), nil
}
