package models

type InputUserSignIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type InputUserSignUp struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	Email               string `json:"email"`
	AdminAuthToken      string `json:"adminAuthToken"`
	RegisterAsModerator bool   `json:"registerAsModerator"`
}

type InputPost struct {
	ID         int64    `json:"id"`
	AuthorID   int64    `json:"authorId"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
	IsImage    bool     `json:"isImage"`
	ImagePath  string   `json:"imagePath"`
	Bans       []string `json:"bans"`
}

type InputComment struct {
	ID       int64  `json:"id"`
	AuthorID int64  `json:"authorId"`
	PostID   int64  `json:"postId"`
	Content  string `json:"content"`
}

type InputFindComment struct {
	Option string `json:"option"` // user or post
	PostID int64  `json:"postId"`
	UserID int64  `json:"userId"`
}

type InputRate struct {
	ID       int64 `json:"id"`       // postId
	Reaction int   `json:"reaction"` // 1 or -1
}

type InputCommentRate struct {
	CommentID int64 `json:"commentId"` // commentId
	PostID    int64 `json:"postId"`    // postId
	Reaction  int   `json:"reaction"`  // 1 or -1
}

type InputFilterPost struct {
	Option     string   `json:"option"`   // categories or date or rating or author or banned
	AuthorID   int64    `json:"authorId"` // getAllPosts created by AuthorId
	Date       string   `json:"date"`     // ASC or DESC
	Rating     string   `json:"rating"`   // ASC or DESC
	Categories []string `json:"categories"`
	UserRating string   `json:"userRating"` // upvoted or downvoted
	UserID     int64    `json:"userId"`
}
