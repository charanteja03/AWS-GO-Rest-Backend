package models

// UserDetailsDto - struct representing user details sent as part of Response
type UserDetailsDto struct {
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Username       string `json:"username"`
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	Email          string `json:"email"`
	TokenExpiresAt string `json:"expiresAt"`
}
