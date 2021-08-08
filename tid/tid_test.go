package tid_test

import (
	"net/http"
	"sfr-backend/tid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTidGeneration(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	expectedTid := "tid"
	req.Header.Add("X-Request-ID", expectedTid)

	tid := tid.GetTid(req)

	assert.Equal(t, expectedTid, tid)
}
