package models

type RoleRequest struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"userID"`
	CreatedAt int64 `json:"createdAt,omitempty"`
	Pending   bool  `json:"pending"`
	User      *User `json:"user"`
}
