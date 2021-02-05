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
}
