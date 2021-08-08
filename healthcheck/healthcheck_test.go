package healthcheck_test

import (
	"net/http"
	"net/http/httptest"
	"sfr-backend/healthcheck"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	rr := httptest.NewRecorder()

	healthcheck.Healthcheck(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
