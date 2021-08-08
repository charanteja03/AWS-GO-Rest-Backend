package models

// RefreshToken -
type RefreshToken struct {
	RefreshToken string `json:"refreshtoken"`
}

// ResponseObject -
type ResponseObject struct {
	RefreshToken    string `json:"refreshToken"`
	AccessToken     string `json:"accessToken"`
	TokenExpiryTime string `json:"expiresAt"`
}
