package models

type Comment struct {
	ID            int64  `json:"id"`
	PostID        int64  `json:"post_id"`
	AuthorID      int64  `json:"-"`
	Content       string `json:"content"`
	CreatedAt     int64  `json:"createdAt"`
	EditedAt      int64  `json:"editedAt"`
	Author        *User  `json:"author"`
	CommentRating int    `json:"commentRating"`
	UserRating    int    `json:"userRating"`
}
