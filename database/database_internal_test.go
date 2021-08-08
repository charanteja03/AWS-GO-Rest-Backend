package database

import (
	"errors"
	"sfr-backend/mocks"
	"sfr-backend/user"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchDatabaseSession(t *testing.T) {
	mockAwsDatabseProvider := &mocks.AwsDatabaseProvider{}
	mockAwsDatabase := &mocks.AwsDatabaseInterface{}
	mockAwsDatabseProvider.On("New", mock.Anything).Return(mockAwsDatabase, nil)
	// change database provider for mock provider
	awsDatabaseProvider = mockAwsDatabseProvider

	databaseProvider := fetchAwsSession()

	assert.Equal(t, mockAwsDatabase, databaseProvider)
}

func TestCreateUser(t *testing.T) {
	mockAwsDatabseProvider := &mocks.AwsDatabaseProvider{}
	mockAwsDatabase := &mocks.AwsDatabaseInterface{}
	mockAwsDatabase.On("PutItem", mock.Anything).Return(nil, errors.New("Error")).Once()
	mockAwsDatabase.On("PutItem", mock.Anything).Return(&dynamodb.PutItemOutput{}, nil).Once()
	mockAwsDatabseProvider.On("New", mock.Anything).Return(mockAwsDatabase, nil)

	awsDatabaseProvider = mockAwsDatabseProvider

	preparedUserDetails := user.UserDetails{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Username:  "Username",
	}
	testTable := []struct {
		userDetails    user.UserDetails
		expectedOutput string
	}{
		{user.UserDetails{}, "error"},
		{preparedUserDetails, "success"},
	}

	for _, testCase := range testTable {
		output := CreateUser(testCase.userDetails)
		assert.Equal(t, testCase.expectedOutput, output)
	}
}

func TestGetUserDetails(t *testing.T) {
	mockAwsDatabseProvider := &mocks.AwsDatabaseProvider{}
	mockAwsDatabase := &mocks.AwsDatabaseInterface{}
	username := "username"
	userReturnValue := make(map[string]*dynamodb.AttributeValue)
	userReturnValue["username"] = &dynamodb.AttributeValue{S: &username}
	userReturnValue["password"] = &dynamodb.AttributeValue{S: &username}
	mockAwsDatabase.On("GetItem", mock.Anything).Return(nil, errors.New("Error")).Once()
	mockAwsDatabase.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{
		Item: userReturnValue,
	}, nil).Once()
	mockAwsDatabase.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{}, nil).Once()
	mockAwsDatabseProvider.On("New", mock.Anything).Return(mockAwsDatabase, nil)

	databaseUser := user.User{}
	dynamodbattribute.UnmarshalMap(userReturnValue, &databaseUser)

	awsDatabaseProvider = mockAwsDatabseProvider

	testTable := []struct {
		expectedUser user.UserDetails
	}{
		{user.UserDetails{}},
		{user.UserDetails{Username: username, Password: username}},
		{user.UserDetails{}},
	}
	for _, testCase := range testTable {
		output := GetUserDetails(username)
		assert.Equal(t, testCase.expectedUser, output)
	}
}
