package docs

import "github.com/aws/aws-sdk-go/service/sfn"

// swagger:route GET /aws/machines machines-endpoint idMachinesEndpoint
// Returns machine's list from current AWS environment.
// responses:
//   200: machinesResponse

// swagger:parameters idMachinesEndpoint
type machinesWrapper struct {
	// Max returned value count.
	// in:query
	// name:count
	// required:false
	Count int32 `json:"count"`
	// Token for AWS pagination.
	// in:query
	// name:nextToken
	// required:false
	NextToken string `json:"nextToken"`
	// JWT Token for authentication in subsequent operations
	// in:header
	Authentication string
}

// Returns a JSON with an element STATUS with the OK text.
// swagger:response machinesResponse
type machinesResponse struct {
	// in:body
	Body sfn.ListStateMachinesOutput
}
