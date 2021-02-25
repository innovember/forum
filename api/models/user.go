package models

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	Email      string `json:"email"`
	CreatedAt  int64  `json:"createdAt,omitempty"`
	LastActive int64  `json:"lastActive,omitempty"`
	SessionID  string `json:"sessionID,omitempty"`
	Role       int    `json:"role"`
}
