package awsprovider

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//go:generate mockery --output=../mocks --name=AwsDatabaseInterface|AwsDatabaseProvider

//AwsDatabaseInterface - interface for aws step functions functionality
type AwsDatabaseInterface interface {
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}

//AwsDatabaseProvider - provider for step function interface
type AwsDatabaseProvider interface {
	New(sess *session.Session) (AwsDatabaseInterface, error)
}

//AwsDatabaseProviderRealProvider - struct represents real AWS provider
type AwsDatabaseProviderRealProvider struct {
}

//New - function creates a new instance of the SFN client with a session.
func (awsDatabaseProvider *AwsDatabaseProviderRealProvider) New(sess *session.Session) (AwsDatabaseInterface, error) {
	return dynamodb.New(sess), nil
}
