package machine_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sfr-backend/machine"
	"sfr-backend/mocks"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProperMachine(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	token := ""
	machineName := "Machine"
	machineArn := "MachineArn"
	machineType := "MachineType"
	stateMachinesList := []*sfn.StateMachineListItem{
		{
			CreationDate:    &time.Time{},
			Name:            &machineName,
			StateMachineArn: &machineArn,
			Type:            &machineType,
		},
	}
	mockOutput := sfn.ListStateMachinesOutput{
		NextToken:     &token,
		StateMachines: stateMachinesList,
	}

	mockStepFunction.On("ListStateMachines", mock.Anything).Return(&mockOutput, nil)

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)

	nextToken := "token"
	count := 2
	testTable := []struct {
		count      *int
		nextToken  *string
		shouldPass bool
	}{
		{&count, nil, true},
		{nil, &nextToken, false},
		{&count, &nextToken, false},
	}
	for _, testCase := range testTable {
		path := fmt.Sprintf("/aws/machines?")
		if testCase.count != nil {
			path += fmt.Sprintf("count=%d", *testCase.count)
		}
		if testCase.nextToken != nil {
			path += fmt.Sprintf("&nextToken=%s", *testCase.nextToken)
		}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			machine.GetMachinesHandler(w, r, mockAwsProvider)
		})
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	}
}

func TestGetErrorInsteadMachine(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}

	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)
	mockStepFunction.On("ListStateMachines", mock.Anything).Return(nil, errors.New("error"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		machine.GetMachinesHandler(w, r, mockAwsProvider)
	})

	path := fmt.Sprintf("/aws/machines?")
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetMachinesStepFunctionProviderError(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}

	mockAwsProvider.On("New", mock.Anything).Return(nil, errors.New("Error"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		machine.GetMachinesHandler(w, r, mockAwsProvider)
	})
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
