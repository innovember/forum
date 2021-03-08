package models

type Notification struct {
	ID            int64          `json:"id"`
	ReceiverID    int64          `json:"receiver_id"`
	PostID        int64          `json:"post_id"`
	RateID        int64          `json:"rate_id"`
	CommentID     int64          `json:"comment_id"`
	CommentRateID int64          `json:"comment_rate_id"`
	CreatedAt     int64          `json:"createdAt"`
	Post          *Post          `json:"post"`
	PostRating    *PostRating    `json:"postRating"`
	Comment       *Comment       `json:"comment"`
	CommentRating *CommentRating `json:"commentRating"`
}

type RoleNotification struct {
	ID         int64 `json:"id"`
	ReceiverID int64 `json:"receiver_id"`
	Accepted   bool  `json:"accepted"`
	Declined   bool  `json:"declined"`
	Demoted    bool  `json:"demoted"`
	CreatedAt  int64 `json:"createdAt,omitempty"`
}

type PostNotification struct {
	ID         int64 `json:"id"`
	ReceiverID int64 `json:"receiver_id"`
	Approved   bool  `json:"approved"`
	Banned     bool  `json:"banned"`
	Deleted    bool  `json:"deleted"`
	CreatedAt  int64 `json:"createdAt,omitempty"`
}
