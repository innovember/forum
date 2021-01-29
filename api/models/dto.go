package models

type InputUserSignIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type InputUserSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type InputPost struct {
	ID         int64    `json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}

type InputComment struct {
	ID      int64  `json:"id"`
	PostID  int64  `json:"post_id"`
	Content string `json:"content"`
}

type InputRate struct {
	ID       int64 `json:"id"`
	Reaction int   `json:"reaction"`
}
