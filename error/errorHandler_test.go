package error_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sfr-backend/error"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	expectedRecorder := httptest.NewRecorder()

	errorMessage := "errorMessage"
	error.HandleError(rr, errors.New(errorMessage))
	http.Error(expectedRecorder, errorMessage, http.StatusBadRequest)

	assert.Equal(t, expectedRecorder.Body, rr.Body)
}
