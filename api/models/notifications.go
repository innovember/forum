package models

type Notification struct {
	ID            int64          `json:"id"`
	ReceiverID    int64          `json:"receiverId"`
	PostID        int64          `json:"postId"`
	RateID        int64          `json:"rateId"`
	CommentID     int64          `json:"commentId"`
	CommentRateID int64          `json:"commentRateId"`
	CreatedAt     int64          `json:"createdAt"`
	Post          *Post          `json:"post"`
	PostRating    *PostRating    `json:"postRating"`
	Comment       *Comment       `json:"comment"`
	CommentRating *CommentRating `json:"commentRating"`
}

type RoleNotification struct {
	ID         int64 `json:"id"`
	ReceiverID int64 `json:"receiverId"`
	Accepted   bool  `json:"accepted"`
	Declined   bool  `json:"declined"`
	Demoted    bool  `json:"demoted"`
	CreatedAt  int64 `json:"createdAt,omitempty"`
}

type PostReportNotification struct {
	ID         int64 `json:"id"`
	ReceiverID int64 `json:"receiverId"`
	Approved   bool  `json:"approved"`
	Deleted    bool  `json:"deleted"`
	CreatedAt  int64 `json:"createdAt,omitempty"`
}

type PostNotification struct {
	ID         int64 `json:"id"`
	ReceiverID int64 `json:"receiverId"`
	Approved   bool  `json:"approved"`
	Banned     bool  `json:"banned"`
	Deleted    bool  `json:"deleted"`
	CreatedAt  int64 `json:"createdAt,omitempty"`
}
