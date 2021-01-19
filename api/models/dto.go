package models

type InputUserSignIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type InputUserSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
