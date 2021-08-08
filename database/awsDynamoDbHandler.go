package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	awsprovider "sfr-backend/awsProvider"
	"sfr-backend/user"
)

var awsDatabaseProvider awsprovider.AwsDatabaseProvider

func init() {
	awsDatabaseProvider = &awsprovider.AwsDatabaseProviderRealProvider{}
}

//GetUserDetails - get user details from DB
func GetUserDetails(username string) user.UserDetails {
	svc := fetchAwsSession()
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("UserDetails"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return user.UserDetails{}
	}

	databaseUser := user.UserDetails{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &databaseUser)

	if err != nil {
		panic(fmt.Sprintf("Failed to Unmarshal JSON Record, %s", err))
	}

	if databaseUser.Password == "" {
		fmt.Println("Could not retrieve user details")
		return user.UserDetails{}
	}

	return databaseUser
}

//CreateUser - function to create a user
func CreateUser(userDetails user.UserDetails) string {
	svc := fetchAwsSession()

	av, error := dynamodbattribute.MarshalMap(userDetails)
	if error != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(error.Error())
		return "error"
	}

	_, err := svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("UserDetails"),
	})

	if err != nil {
		fmt.Println(err.Error())
		return "error"
	}

	fmt.Println("Successfully added user with username " + userDetails.Username + " to table UserDetails")

	return "success"
}

func fetchAwsSession() awsprovider.AwsDatabaseInterface {

	sess, err := session.NewSessionWithOptions(session.Options{
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to Unmarshal JSON Record, %s", err))
	}
	svc, _ := awsDatabaseProvider.New(sess)

	return svc
}
