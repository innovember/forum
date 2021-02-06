package models

type Post struct {
	ID             int64      `json:"id"`
	AuthorID       int64      `json:"-"`
	Author         *User      `json:"author"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Categories     []Category `json:"categories"`
	PostRating     int        `json:"postRating"`
	UserRating     int        `json:"userRating"`
	CreatedAt      int64      `json:"createdAt,omitempty"`
	CommentsNumber int        `json:"commentsNumber"`
}
