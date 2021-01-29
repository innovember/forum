package models

type Category struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PostAttached int64  `json:"postAttached,omitempty"`
}
