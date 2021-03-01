package models

type PostReport struct {
	ID          int64 `json:"id"`
	ModeratorID int64 `json:"moderatorID"`
	PostID      int64 `json:"postID"`
	CreatedAt   int64 `json:"createdAt,omitempty"`
	Pending     bool  `json:"pending"`
}
