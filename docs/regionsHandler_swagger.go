package docs

// swagger:route GET /aws/regions regions-endpoint idGetRegionsEndpoint
// Returns the complete list of Regions.
// responses:
//   200: getRegionsResponse

// swagger:parameters idGetRegionsEndpoint
type getRegionsWrapper struct {
	// JWT value for authentication in subsequent operations
	// in:header
	Authentication string
}

// Returns region's list.
// swagger:response getRegionsResponse
type getRegionsResponse struct {
	// in:body
	region []string
}
