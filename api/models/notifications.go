package models

type Notification struct {
	ID         int64 `json:"id"`
	ReceiverID int64 `json:"receiver_id"`
	PostID     int64 `json:"post_id"`
	RateID     int64 `json:"rate_id"`
	CommentID  int64 `json:"comment_id"`
	CreatedAt  int64 `json:"createdAt"`
}
