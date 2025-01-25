package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hunderaweke/metsasft/internal/domain"
)

type BlogController struct {
	usecase domain.BlogUsecase
}

func NewBlogController(usecase domain.BlogUsecase) *BlogController {
	return &BlogController{usecase: usecase}
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var blog domain.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(blog); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	blog.WriterID = userId
	blog, err := c.usecase.CreateBlog(blog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, blog)
}

func (c *BlogController) GetBlogs(ctx *gin.Context) {
	var blogs []domain.Blog
	blogs, err := c.usecase.GetBlogs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}
func (c *BlogController) GetBlog(ctx *gin.Context) {
	id := ctx.Param("blog_id")
	blog, err := c.usecase.GetBlogByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blog)
}
func (c *BlogController) UpdateBlog(ctx *gin.Context) {
	id := ctx.Param("blog_id")
	var blog domain.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(blog); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	blog.ID = id
	blog, err := c.usecase.UpdateBlog(blog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blog)
}
func (c *BlogController) DeleteBlog(ctx *gin.Context) {
	id := ctx.Param("blog_id")
	err := c.usecase.DeleteBlog(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}

func (c *BlogController) GetBlogsByWriterID(ctx *gin.Context) {
	writerID := ctx.Param("writer_id")
	blogs, err := c.usecase.GetBlogByWriterID(writerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}

func (c *BlogController) GetBlogsByStatus(ctx *gin.Context) {
	status := ctx.Param("status")
	blogs, err := c.usecase.GetBlogByStatus(status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}

func (c *BlogController) GetBlogsByModificationDateRange(ctx *gin.Context) {
	var modificationRange domain.ModificationDateRange
	if err := ctx.ShouldBindJSON(&modificationRange); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	blogs, err := c.usecase.GetBlogByModificationDateRange(modificationRange)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}

func (c *BlogController) GetComments(ctx *gin.Context) {
	blogID := ctx.Param("blog_id")
	comments, err := c.usecase.GetComments(blogID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, comments)
}

func (c *BlogController) CreateComment(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	blogID := ctx.Param("blog_id")
	var comment domain.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(comment); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	comment.WriterID = userID
	comment.BlogID = blogID
	comment, err := c.usecase.CreateComment(comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, comment)
}

func (c *BlogController) GetComment(ctx *gin.Context) {
	blogId := ctx.Param("blog_id")
	commentId := ctx.Param("id")
	comment, err := c.usecase.GetCommentByID(blogId, commentId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, comment)
}
