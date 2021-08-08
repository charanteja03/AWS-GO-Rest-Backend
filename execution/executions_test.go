package execution_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sfr-backend/execution"
	"sfr-backend/mocks"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetExecutionsList(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)

	executionName := "executionName"
	executionArn := "executionArn"
	machineArn := "machineArn"
	status := "SUCCEEDED"
	executionsList := []*sfn.ExecutionListItem{
		{
			Name:            &executionName,
			StateMachineArn: &machineArn,
			StartDate:       &time.Time{},
			ExecutionArn:    &executionArn,
			Status:          &status,
		},
	}
	output := &sfn.ListExecutionsOutput{Executions: executionsList}
	mockStepFunction.On("ListExecutions", mock.Anything).Return(output, nil)

	count := 3
	nextToken := "token"
	statusFilter := "FAILED" // possible filters: ["RUNNING", "SUCCEEDED", "FAILED", "TIMED_OUT", "ABORTED"];
	testTable := []struct {
		machine      string
		count        *int
		nextToken    *string
		statusFilter *string
	}{
		{"machine", nil, nil, nil},
		{"machine", &count, &nextToken, nil},
		{"machine", &count, nil, &statusFilter},
	}

	for _, testCase := range testTable {
		path := fmt.Sprintf("/aws/executions?machine=%s", testCase.machine)
		if testCase.count != nil {
			path += fmt.Sprintf("&count=%d", *testCase.count)
		}
		if testCase.nextToken != nil {
			path += fmt.Sprintf("&nextToken=%s", *testCase.nextToken)
		}
		if testCase.statusFilter != nil {
			path += fmt.Sprintf("&statusFilter=%s", *testCase.statusFilter)
		}
		req, _ := http.NewRequest("GET", path, nil)
		// We create a ResponseRecorder to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			execution.GetExecutionsHandler(w, r, mockAwsProvider)
		})
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	}
}

func TestGetExecutionsListSessionCreationError(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}

	errorMessage := "Error"
	mockAwsProvider.On("New", mock.Anything).Return(nil, errors.New(errorMessage))

	path := fmt.Sprintf("/aws/executions?machine=%s", "machine")
	req, _ := http.NewRequest("GET", path, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.GetExecutionsHandler(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetExecutionsListError(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)
	errorMessage := "Error"
	mockStepFunction.On("ListExecutions", mock.Anything).Return(nil, errors.New(errorMessage))

	path := fmt.Sprintf("/aws/executions?machine=%s", "machine")
	req, _ := http.NewRequest("GET", path, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.GetExecutionsHandler(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetExecutionHandler(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)

	executionArn := "executionArn"
	output := &sfn.DescribeExecutionOutput{
		ExecutionArn: &executionArn,
	}
	mockStepFunction.On("DescribeExecution", mock.Anything).Return(output, nil)

	path := fmt.Sprintf("/aws/execution/%s", "execution")
	req, _ := http.NewRequest("GET", path, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.GetExecutionHandler(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetExecutionHandlerSessionCreationError(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}

	errorMessage := "Error"
	mockAwsProvider.On("New", mock.Anything).Return(nil, errors.New(errorMessage))

	path := fmt.Sprintf("/aws/execution/%s", "execution")
	req, _ := http.NewRequest("GET", path, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.GetExecutionHandler(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetExecutionHandlerProviderError(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)

	errorMessage := "errorMessage"
	mockStepFunction.On("DescribeExecution", mock.Anything).Return(nil, errors.New(errorMessage))

	path := fmt.Sprintf("/aws/execution/%s", "execution")
	req, _ := http.NewRequest("GET", path, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.GetExecutionHandler(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPostStartExecutionNoInput(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)
	executionArn := "executionArn"
	output := &sfn.StartExecutionOutput{
		ExecutionArn: &executionArn,
		StartDate:    &time.Time{},
	}
	mockStepFunction.On("StartExecution", mock.Anything).Return(output, nil)

	payload := strings.NewReader("machine=machine")
	req, _ := http.NewRequest("POST", "/aws/execution", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.PostStartExecution(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestPostStartExecutionWithInput(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)
	executionArn := "executionArn"
	output := &sfn.StartExecutionOutput{
		ExecutionArn: &executionArn,
		StartDate:    &time.Time{},
	}
	mockStepFunction.On("StartExecution", mock.Anything).Return(output, nil)

	payload := strings.NewReader("machine=machine&input=input")
	req, _ := http.NewRequest("POST", "/aws/execution", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.PostStartExecution(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestPostStartExecutionSessionCreationError(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}

	errorMessage := "errorMessage"
	mockAwsProvider.On("New", mock.Anything).Return(nil, errors.New(errorMessage))

	payload := strings.NewReader("machine=machine")
	req, _ := http.NewRequest("POST", "/aws/execution", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.PostStartExecution(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPostStartExecutionProviderError(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)
	errorMessage := "errorMessage"
	mockStepFunction.On("StartExecution", mock.Anything).Return(nil, errors.New(errorMessage))

	payload := strings.NewReader("machine=machine")
	req, _ := http.NewRequest("POST", "/aws/execution", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.PostStartExecution(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPostRestartExecution(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)

	executionArn := "executionArn"
	input := "input"
	describeExecutionoutput := &sfn.DescribeExecutionOutput{
		ExecutionArn: &executionArn,
		Input:        &input,
	}
	mockStepFunction.On("DescribeExecution", mock.Anything).Return(describeExecutionoutput, nil)

	startExecutionOutput := &sfn.StartExecutionOutput{
		ExecutionArn: &executionArn,
		StartDate:    &time.Time{},
	}
	mockStepFunction.On("StartExecution", mock.Anything).Return(startExecutionOutput, nil)

	payload := strings.NewReader("machine=machine&execution=execution")
	req, _ := http.NewRequest("POST", "/aws/execution", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.PostRestartExecution(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestPostRestartExecutionProviderError(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)
	errorMessage := "errorMessage"
	mockStepFunction.On("DescribeExecution", mock.Anything).Return(nil, errors.New(errorMessage))

	payload := strings.NewReader("machine=machine&execution=execution")
	req, _ := http.NewRequest("POST", "/aws/execution", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.PostRestartExecution(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPostRestartBatch(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)
	errorMessage := "errorMessage"
	mockStepFunction.On("StartExecution", mock.Anything).Return(nil, errors.New(errorMessage))

	payload := strings.NewReader("machine=machine&executions=[\"execution\"]&useOriginalInput=false&input={}")
	req, _ := http.NewRequest("POST", "/aws/execution", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execution.PostRestartBatch(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
