package execution

import (
	"encoding/json"
	"net/http"
	"strconv"

	awsprovider "sfr-backend/awsProvider"
	"sfr-backend/awssession"
	errHandler "sfr-backend/error"
	"sfr-backend/response"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/gorilla/mux"
)

// GetExecutionsHandler - returns all executions on given machine filtered with statusFilter
func GetExecutionsHandler(w http.ResponseWriter, r *http.Request, providerInterface awsprovider.AwsStepFunctionsProvider) {
	vars := mux.Vars(r)

	sfv, err := awssession.CreateStepFunctionSession(w, r, providerInterface)
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}

	input := &sfn.ListExecutionsInput{
		StateMachineArn: aws.String(vars["machine"]),
	}
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

	if len(urlParams.Get("statusFilter")) > 0 {
		statusToSet := urlParams.Get("statusFilter")
		executionsStates := sfn.ExecutionStatus_Values()
		correctStatus := false
		for _, status := range executionsStates {
			if status == statusToSet {
				correctStatus = true
				break
			}
		}
		if correctStatus {
			input.StatusFilter = aws.String(statusToSet)
		}
	}
	executions, err := sfv.ListExecutions(input)
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}
	// fmt.Println("executions", executions)
	response.WriteResponse(w, executions)
}

// GetExecutionHandler - returns execution details containing execution input and output
func GetExecutionHandler(w http.ResponseWriter, r *http.Request, providerInterface awsprovider.AwsStepFunctionsProvider) {
	vars := mux.Vars(r)

	sfv, err := awssession.CreateStepFunctionSession(w, r, providerInterface)
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}

	executions, err := sfv.DescribeExecution(&sfn.DescribeExecutionInput{
		ExecutionArn: aws.String(vars["execution"]),
	})

	if err != nil {
		errHandler.HandleError(w, err)
		return
	}

	// fmt.Println("executions", executions)
	response.WriteResponse(w, executions)
}

// PostStartExecution - starts executions with given params
func PostStartExecution(w http.ResponseWriter, r *http.Request, providerInterface awsprovider.AwsStepFunctionsProvider) {
	sfv, err := awssession.CreateStepFunctionSession(w, r, providerInterface)
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}
	err = r.ParseForm()

	executionInput := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(r.FormValue("machine")),
	}

	if len(r.FormValue("input")) > 0 {
		executionInput.Input = aws.String(r.FormValue("input"))
	} else {
		executionInput.Input = aws.String("{}")
	}
	executionStart, err := sfv.StartExecution(executionInput)

	if err != nil {
		errHandler.HandleError(w, err)
		return
	}
	response.WriteResponse(w, executionStart)
}

// PostRestartExecution - restarts given execution
func PostRestartExecution(w http.ResponseWriter, r *http.Request, providerInterface awsprovider.AwsStepFunctionsProvider) {
	sfv, err := awssession.CreateStepFunctionSession(w, r, providerInterface)
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}
	err = r.ParseForm()
	executionStart, err := runExecution(sfv, r.FormValue("machine"), r.FormValue("execution"), "")
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}
	response.WriteResponse(w, executionStart)
}

// PostRestartBatch - post request to reproces execution batch
func PostRestartBatch(w http.ResponseWriter, r *http.Request, providerInterface awsprovider.AwsStepFunctionsProvider) {
	sfv, err := awssession.CreateStepFunctionSession(w, r, providerInterface)
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}
	err = r.ParseForm()
	useOriginalInput, err := strconv.ParseBool(r.FormValue("useOriginalInput"))
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}
	var executions []string
	err = json.Unmarshal([]byte(r.FormValue("executions")), &executions)
	if err != nil {
		errHandler.HandleError(w, err)
		return
	}
	executionsBatches := []*sfn.StartExecutionOutput{}
	errors := []string{}
	for _, execution := range executions {
		if useOriginalInput {
			rerun, err := runExecution(sfv, r.FormValue("machine"), execution, "")
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				executionsBatches = append(executionsBatches, rerun)
			}
		} else {
			rerun, err := runExecution(sfv, r.FormValue("machine"), execution, r.FormValue("input"))
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				executionsBatches = append(executionsBatches, rerun)
			}
		}
	}
	type ResponseStruct struct {
		Execution []*sfn.StartExecutionOutput
		Errors    []string
	}
	responseData := ResponseStruct{
		Execution: executionsBatches,
		Errors:    errors,
	}
	response.WriteResponse(w, responseData)
}

func runExecution(stepFunctionAPI awsprovider.AwsStepFunctionInterface, machine string, execution string, input string) (*sfn.StartExecutionOutput, error) {
	executionInput := &sfn.StartExecutionInput{
		StateMachineArn: aws.String(machine),
	}
	// if input is provided use input for execution
	// else get last input from execution and use as input
	if len(input) > 0 {
		executionInput.Input = aws.String(input)
	} else {
		// get execution input
		execution, err := stepFunctionAPI.DescribeExecution(&sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(execution),
		})
		if err != nil {
			return nil, err
		}
		executionInput.Input = aws.String(*execution.Input)
	}
	executionStart, err := stepFunctionAPI.StartExecution(executionInput)
	if err != nil {
		return nil, err
	}
	return executionStart, nil

}
