package user

// User = Used to store User Details
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserDetails - struct representing user details
type UserDetails struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}
