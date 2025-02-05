package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/metsasft/api/controllers"
	"github.com/hunderaweke/metsasft/api/middlewares"
	"github.com/hunderaweke/metsasft/internal/usecases"
	"github.com/sv-tools/mongoifc"
)

func AddBlogRoutes(r *gin.Engine, db mongoifc.Database, ctx context.Context) {
	usecase := usecases.NewBlogUsecase(db, ctx)
	controller := controllers.NewBlogController(usecase)
	blogRouter := r.Group("/blogs")
	{
		blogRouter.GET("/", controller.GetBlogs)
		blogRouter.GET("/:blog_id", controller.GetBlog)
		blogRouter.GET("/:blog_id/comments", controller.GetComments)
		blogRouter.GET("/:blog_id/comments/:comment_id", controller.GetComment)
	}
	blogRouter.Use(middlewares.AuthenticationMiddleware())
	{
		blogRouter.POST("/", controller.CreateBlog)
		blogRouter.PUT("/:blog_id", controller.UpdateBlog)
		blogRouter.POST("/:blog_id/comments", controller.CreateComment)
		blogRouter.DELETE("/:blog_id", controller.DeleteBlog, middlewares.IsAdmin())
	}
}
