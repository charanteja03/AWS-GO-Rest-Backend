package awssession_test

import (
	"net/http"
	"net/http/httptest"
	"sfr-backend/awssession"
	"sfr-backend/mocks"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var awsSessionMock func(sess *session.Session) *sfn.SFN

type preSessionMock struct{}

func (u preSessionMock) New(sess *session.Session) *sfn.SFN {
	return awsSessionMock(sess)
}

func TestStepFunctionSession(t *testing.T) {
	mockAwsProvider := &mocks.AwsStepFunctionsProvider{}
	mockStepFunction := &mocks.AwsStepFunctionInterface{}
	mockAwsProvider.On("New", mock.Anything).Return(mockStepFunction, nil)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	stepFunctionProvider, error := awssession.CreateStepFunctionSession(rr, req, mockAwsProvider)

	assert.Equal(t, mockStepFunction, stepFunctionProvider)
	assert.Equal(t, error, nil)
}
