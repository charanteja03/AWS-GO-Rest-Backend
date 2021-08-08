package response_test

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http/httptest"
	"sfr-backend/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {

	rr := httptest.NewRecorder()
	result := "result"
	expectedOutput, _ := json.Marshal(result)
	response.WriteResponse(rr, result)

	assert.Equal(t, string(expectedOutput[:]), rr.Body.String())
}

func TestWriteResponseWithError(t *testing.T) {

	rr := httptest.NewRecorder()
	errorProne := math.Inf(1)
	_, error := json.Marshal(errorProne)
	response.WriteResponse(rr, errorProne)

	var b bytes.Buffer
	b.Write([]byte(error.Error()))
	//Adding new line becouse buffer is flushing in WriteResponse
	b.WriteString("\n")

	assert.Equal(t, &b, rr.Body)
}
