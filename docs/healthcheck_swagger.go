package docs

import (
	"sfr-backend/healthcheck"
)

// swagger:route GET /healthcheck healthcheck-endpoint healthCheckEndpoint
// Healthcheck returns a fixed JSON structure for healthcheck.
// responses:
//   200: healthcheckResponse

// Returns a JSON with an element STATUS with the OK text.
// swagger:response healthcheckResponse
type healthcheckResponse struct {
	// in:body
	Body healthcheck.ResponseHealthcheck
	// Set-Cookie value for authentication in subsequent operations
	Cookie string
}
