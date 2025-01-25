package domain

import "time"

const BlogCollection = "blogs"

type Blog struct {
	ID               string    `bson:"_id,omitempty" json:"id"`
	Title            string    `json:"title" validate:"required"`
	Body             string    `json:"body" validate:"required"`
	WriterID         string    `bson:"writer_id" json:"writer_id" validate:"required"`
	Status           string    `json:"status" validate:"required"`
	LastModifiedDate time.Time `bson:"last_modified_date" json:"last_modified_date" validate:"required"`
}
type ModificationDateRange struct {
	StartDate time.Time `bson:"$gte" json:"start_date"`
	EndDate   time.Time `bson:"$lte" json:"end_date"`
}

type BlogFilter struct {
	Status                string                `bson:"status" json:"status"`
	ModificationDateRange ModificationDateRange `bson:"last_modified_date" json:"last_modified_date"`
	WriterID              string                `bson:"writer_id" json:"writer_id"`
}

type BlogRepository interface {
	CreateBlog(blog Blog) (Blog, error)
	GetBlogs(filter BlogFilter) ([]Blog, error)
	UpdateBlog(blog Blog) (Blog, error)
	DeleteBlog(id string) error
	GetBlogByID(id string) (Blog, error)
}

type BlogUsecase interface {
	CreateBlog(blog Blog) (Blog, error)
	GetBlogs() ([]Blog, error)
	UpdateBlog(blog Blog) (Blog, error)
	DeleteBlog(id string) error
	GetBlogByID(id string) (Blog, error)
	GetBlogByWriterID(writerID string) ([]Blog, error)
	GetBlogByStatus(status string) ([]Blog, error)
	GetBlogByModificationDateRange(ModificationDateRange) ([]Blog, error)
	GetComments(blogID string) ([]Comment, error)
	CreateComment(comment Comment) (Comment, error)
	GetCommentByID(id string) (Comment, error)
}
