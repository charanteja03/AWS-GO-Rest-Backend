package machine

import (
	"net/http"
	awsprovider "sfr-backend/awsProvider"
	"sfr-backend/awssession"
	"sfr-backend/error"
	"sfr-backend/response"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
)

// GetMachinesHandler - returns list of step machines
func GetMachinesHandler(w http.ResponseWriter, r *http.Request, providerInterface awsprovider.AwsStepFunctionsProvider) {
	sfv, err := awssession.CreateStepFunctionSession(w, r, providerInterface)
	if err != nil {
		error.HandleError(w, err)
		return
	}
	input := &sfn.ListStateMachinesInput{}
	urlParams := r.URL.Query()
	if len(urlParams.Get("count")) > 0 {
		countToParse := urlParams.Get("count")
		count, err := strconv.ParseInt(countToParse, 10, 64)
		if err == nil {
			input.MaxResults = &count
		}
	}
	if len(urlParams.Get("nextToken")) > 0 {
		input.NextToken = aws.String(urlParams.Get("nextToken"))
	}
	machines, err := sfv.ListStateMachines(input)
	if err != nil {
		error.HandleError(w, err)
		return
	}

	// fmt.Println("machines", machines)
	response.WriteResponse(w, machines)
}
