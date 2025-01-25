package usecases

import (
	"context"

	"github.com/hunderaweke/metsasft/internal/domain"
	"github.com/hunderaweke/metsasft/internal/repository"
	"github.com/sv-tools/mongoifc"
)

type blogUsecase struct {
	blogRepo    domain.BlogRepository
	commentRepo domain.CommentRepository
}

func NewBlogUsecase(db mongoifc.Database, ctx context.Context) domain.BlogUsecase {
	blogRepo := repository.NewBlogRepository(db, ctx)
	commentRepo := repository.NewCommentRepository(db, ctx)
	return &blogUsecase{blogRepo: blogRepo, commentRepo: commentRepo}
}
func (u *blogUsecase) CreateBlog(blog domain.Blog) (domain.Blog, error) {
	return u.blogRepo.CreateBlog(blog)
}
func (u *blogUsecase) GetBlogs() ([]domain.Blog, error) {
	return u.blogRepo.GetBlogs(domain.BlogFilter{})
}
func (u *blogUsecase) UpdateBlog(blog domain.Blog) (domain.Blog, error) {
	return u.blogRepo.UpdateBlog(blog)
}
func (u *blogUsecase) DeleteBlog(id string) error {
	return u.blogRepo.DeleteBlog(id)
}
func (u *blogUsecase) GetBlogByID(id string) (domain.Blog, error) {
	return u.blogRepo.GetBlogByID(id)
}
func (u *blogUsecase) GetBlogByWriterID(writerID string) ([]domain.Blog, error) {
	return u.blogRepo.GetBlogs(domain.BlogFilter{WriterID: writerID})
}
func (u *blogUsecase) GetBlogByStatus(status string) ([]domain.Blog, error) {
	return u.blogRepo.GetBlogs(domain.BlogFilter{Status: status})
}
func (u *blogUsecase) GetBlogByModificationDateRange(modificationRange domain.ModificationDateRange) ([]domain.Blog, error) {
	return u.blogRepo.GetBlogs(domain.BlogFilter{ModificationDateRange: modificationRange})
}
func (u *blogUsecase) GetComments(blogID string) ([]domain.Comment, error) {
	return u.commentRepo.GetComments(domain.CommentFilter{BlogID: blogID})
}
func (u *blogUsecase) CreateComment(comment domain.Comment) (domain.Comment, error) {
	return u.commentRepo.CreateComment(comment)
}
func (u *blogUsecase) GetCommentByID(blogId, commentId string) (domain.Comment, error) {
	return u.commentRepo.GetCommentByID(blogId, commentId)
}
