package models

type Rating struct {
	Rating     int `json:"rating"`
	UserRating int `json:"userRating"`
}

type PostRating struct {
	ID     int   `json:"id"`
	UserID int64 `json:"userID"`
	PostID int64 `json:"postID"`
	Rate   int   `json:"rate"`
	Author *User `json:"author"`
}

type CommentRating struct {
	ID        int   `json:"id"`
	UserID    int64 `json:"userID"`
	CommentID int64 `json:"commentID"`
	Rate      int   `json:"rate"`
	PostID    int64 `json:"postID"`
	Author    *User `json:"author"`
}
