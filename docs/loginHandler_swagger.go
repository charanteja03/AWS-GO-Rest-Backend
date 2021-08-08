package docs

import (
	"sfr-backend/models"
	"sfr-backend/user"
)

// swagger:route POST /login users-endpoint idLoginEndpoint
// login returns a fixed string SUCCESS when user is properly logged in.
// responses:
//   200: loginResponse

// swagger:parameters idLoginEndpoint
type loginParamsWrapper struct {
	// User structure for login.
	// in:body
	Body user.User
	// Set-Cookie value for authentication in subsequent operations
	Cookie string
}

// Returns a fixed string with a STATUS string
// swagger:response loginResponse
type loginResponse struct {
	// in:body
	result string
	// in: header
}

// Returns a fixed string with a STATUS string indicating error with code 401
// swagger:response authfailureResponse
type authfailureResponse struct {
	result string
}

// swagger:route POST /createuser users-endpoint idCreateUserEndpoint
// Returns a complete user when it is created successfully.
// responses:
//   200: createuserResponse

// swagger:parameters idCreateUserEndpoint
type createUserParamsWrapper struct {
	// Complete User body needed for new user creation.
	// in:body
	Body user.UserDetails
}

// Returns a fixed string with a STATUS string
// swagger:response createuserResponse
type createuserResponse struct {
	result user.UserDetails
}

// swagger:route POST /logout users-endpoint idLogoutEndpoint
// Returns a complete user when it is created successfully.
// responses:
//   200: logoutResponse

// swagger:parameters idLogoutEndpoint
type logoutWrapper struct {
	// JWT value for authentication in subsequent operations
	// in:header
	Authentication string
}

// Returns a fixed string with a message indicating successful logout
// swagger:response logoutResponse
type logoutResponse struct {
	result string
}

// swagger:route POST /refreshtoken users-endpoint idRefreshToken
// Refreshes JWT Token data.
// responses:
//   200: refreshTokenResponse

// swagger:parameters idRefreshToken
type refreshTokenWrapper struct {
	// JWT value for authentication in subsequent operations
	// in:body
	Body models.RefreshToken
}

// Refreshed JWT Token.
// swagger:response refreshTokenResponse
type refreshTokenResponse struct {
	result models.ResponseObject
}
