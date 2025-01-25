package domain

import "time"

const CommentCollection = "comments"

type Comment struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	Body          string    `json:"body" validate:"required"`
	WriterID      string    `json:"writer_id" validate:"required"`
	BlogID        string    `json:"blog_id" validate:"required"`
	CommentedDate time.Time `bson:"commented_date" json:"commented_date" validate:"required"`
}

type CommentFilter struct {
	WriterID string `bson:"writer_id" json:"writer_id"`
	BlogID   string `bson:"blog_id" json:"blog_id"`
}

type CommentRepository interface {
	CreateComment(comment Comment) (Comment, error)
	GetComments(CommentFilter) ([]Comment, error)
	UpdateComment(comment Comment) (Comment, error)
	DeleteComment(id string) error
	GetCommentByID(blogId, commentId string) (Comment, error)
}
