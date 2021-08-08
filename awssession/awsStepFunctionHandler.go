package awssession

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	awsprovider "sfr-backend/awsProvider"
	"sfr-backend/region"
)

// CreateStepFunctionSession - creates session for executiong stepfunctions calls
func CreateStepFunctionSession(w http.ResponseWriter, r *http.Request, awsInterface awsprovider.AwsStepFunctionsProvider) (awsprovider.AwsStepFunctionInterface, error) {
	//Setting some default region for convience
	region := region.GetDefaultRegion(r)

	// an example API handler
	sess, err := session.NewSessionWithOptions(session.Options{
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		return nil, err
	}
	return awsInterface.New(sess)
}
