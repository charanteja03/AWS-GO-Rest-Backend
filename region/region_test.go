package region_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sfr-backend/region"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	godotenv.Load()
}

func TestGetRegionsHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/aws/regions", nil)
	rr := httptest.NewRecorder()

	region.GetRegionsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code to be status code ok")

	var regions []string
	json.Unmarshal(rr.Body.Bytes(), &regions)
	if len(regions) != 26 {
		t.Error("Expected 26 but got ", len(regions))
		assert.Equal(t, 26, len(regions), "Expected to get 26 regions")
	}
}

func TestGetDefaultRegion(t *testing.T) {
	testTable := []struct {
		inputRegion    string
		expectedOutput string
	}{
		{"abc", "abc"},
		{"us-east-1", "us-east-1"},
		{"us-west-1", "us-west-1"},
		{"", "us-east-1"},
	}

	for _, testCase := range testTable {
		path := fmt.Sprintf("/aws?region=%s", testCase.inputRegion)
		// Create a request to pass to our handler.
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}
		returnedRegion := region.GetDefaultRegion(req)
		assert.Equal(t, testCase.expectedOutput, returnedRegion)
	}
}
