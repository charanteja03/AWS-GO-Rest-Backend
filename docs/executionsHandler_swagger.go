package docs

import (
	"github.com/aws/aws-sdk-go/service/sfn"
)

// swagger:route GET /aws/executions executions-endpoint idGetExecutions
// Returns the list of executions for a provided stepfunction.
// responses:
//   200: getExecutionsResponse

// swagger:parameters idGetExecutions
type getExecutionsWrapper struct {
	// State Machine's id to get stepfunction.
	// in:query
	// name:machine
	// required:true
	Machine string `json:"machine"`
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
	// Status value to filter query.
	// in:query
	// name:statusFilter
	// required:false
	StatusFilter string `json:"statusFilter"`
	// JWT Token for authentication in subsequent operations
	// in:header
	Authentication string
}

// Returns a JSON with a list of executions
// swagger:response getExecutionsResponse
type getExecutionsResponse struct {
	// in:body
	Body sfn.ListExecutionsOutput
}

// swagger:route GET /aws/execution/{execution} executions-endpoint idGetExecution
// Returns a specific execution.
// responses:
//   200: getExecutionResponse

// swagger:parameters idGetExecution
type getExecutionWrapper struct {
	// Execution id to get information from.
	// in:path
	// name:execution
	// required:true
	Execution string `json:"execution"`
	// JWT Token for authentication in subsequent operations
	// in:header
	Authentication string
}

// Returns a JSON with the execution output.
// swagger:response getExecutionResponse
type getExecutionResponse struct {
	// in:body
	Body sfn.DescribeExecutionOutput
}

// swagger:route POST /aws/execution executions-endpoint idCreateExecution
// Executes a specific stepfunction with given parameters.
// responses:
//   200: executionResponse

// swagger:parameters idCreateExecution
type executionWrapper struct {
	// Execution id to get information from.
	// in:formData
	// name:machine
	// required:traue
	Machine string `json:"machine"`
	// JSON Formatted Input.
	// in:formData
	// name:input
	// required:true
	Input string `json:"input"`
	// JWT Token for authentication in subsequent operations
	// in:header
	Authentication string
}

// Returns a JSON with the start date and the Execution ARN.
// swagger:response executionResponse
type executionResponse struct {
	// in:body
	Body sfn.StartExecutionOutput
}

// swagger:route POST /aws/execution/restart executions-endpoint idRecreateExecution
// Executes a specific stepfunction with same original parameters.
// responses:
//   200: executionRestartResponse

// swagger:parameters idRecreateExecution
type reexecutionWrapper struct {
	// Execution id to get information from.
	// in:formData
	// name:machine
	// required:traue
	Machine string `json:"machine"`
	// Execution ID for current machine.
	// in:formData
	// name:execution
	// required:true
	Execution string `json:"execution"`
	// JWT Token for authentication in subsequent operations
	// in:header
	Authentication string
}

// Returns a JSON with the start date and the Execution ARN.
// swagger:response executionRestartResponse
type executionRestartResponse struct {
	// in:body
	Body sfn.StartExecutionOutput
}

// swagger:route POST /aws/execution/batch executions-endpoint idBatchExecution
// Rexecutes a list of stepfunctions with original parameters.
// responses:
//   200: executionBatchResponse

// swagger:parameters idBatchExecution
type executionBatchWrapper struct {
	// Execution id to get information from.
	// in:formData
	// name:machine
	// required:traue
	Machine string `json:"machine"`
	// Execution ID for current machine.
	// in:formData
	// name:execution
	// required:true
	Executions []string `json:"executions"`
	// Indicates whether original input should be used or not.
	// in:formData
	// name:useOriginalInput
	// required:true
	UseOriginalInput bool `json:"useOriginalInput"`
	// JSON Formatted Input for all executions.
	// in:formData
	// name:input
	// required:false
	Input string `json:"input"`
	// JWT Token for authentication in subsequent operations
	// in:header
	Authentication string
}

// Returns a JSON list with the start dates and the Executions ARN for all batch executions.
// swagger:response executionBatchResponse
type executionBatchResponse struct {
	// in:body
	Execution []sfn.StartExecutionOutput
	Errors    []string
}
