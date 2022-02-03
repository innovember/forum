package models

type Rating struct {
	Rating     int `json:"rating"`
	UserRating int `json:"userRating"`
}

type PostRating struct {
	ID     int   `json:"id"`
	UserID int64 `json:"userId"`
	PostID int64 `json:"postId"`
	Rate   int   `json:"rate"`
	Author *User `json:"author"`
}

type CommentRating struct {
	ID        int   `json:"id"`
	UserID    int64 `json:"userId"`
	CommentID int64 `json:"commentId"`
	Rate      int   `json:"rate"`
	PostID    int64 `json:"postId"`
	Author    *User `json:"author"`
}
