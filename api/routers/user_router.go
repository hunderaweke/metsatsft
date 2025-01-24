package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/metsasft/api/controllers"
	"github.com/hunderaweke/metsasft/api/middlewares"
	"github.com/hunderaweke/metsasft/internal/usecases"
	"github.com/sv-tools/mongoifc"
)

func AddUserRoutes(r *gin.Engine, db mongoifc.Database, ctx context.Context) error {
	usecase := usecases.NewUserUsecase(db, ctx)
	controller := controllers.NewUserController(usecase)
	r.POST("/login", controller.Login)
	r.POST("/refresh", controller.RefreshToken)
	r.POST("/users", controller.CreateUser)
	userRoutes := r.Group("/users")
	userRoutes.Use(middlewares.AuthenticationMiddleware())
	{
		userRoutes.GET("/", controller.GetUsers)
		userRoutes.GET("/:id", controller.GetUser)
		userRoutes.PUT("/:id", controller.UpdateUser)
		userRoutes.DELETE("/:id", controller.DeleteUser, middlewares.IsAdmin())
		userRoutes.POST("/:id/activate", controller.ActivateUser, middlewares.IsAdmin())
		userRoutes.POST("/:id/deactivate", controller.DeactivateUser, middlewares.IsAdmin())
	}
	return nil
}
