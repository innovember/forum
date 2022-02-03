package models

type PostReport struct {
	ID          int64  `json:"id"`
	ModeratorID int64  `json:"moderatorId"`
	PostID      int64  `json:"postId"`
	CreatedAt   int64  `json:"createdAt,omitempty"`
	Pending     bool   `json:"pending"`
	PostTitle   string `json:"postTitle"`
}
