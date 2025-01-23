package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/metsasft/api/controllers"
	"github.com/hunderaweke/metsasft/internal/usecases"
	"github.com/sv-tools/mongoifc"
)

func AddUserRoutes(r *gin.Engine, db mongoifc.Database, ctx context.Context) error {
	usecase := usecases.NewUserUsecase(db, ctx)
	controller := controllers.NewUserController(usecase)
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", controller.CreateUser)
		userRoutes.GET("/", controller.GetUsers)
		userRoutes.GET("/:id", controller.GetUser)
		userRoutes.PUT("/:id", controller.UpdateUser)
		userRoutes.DELETE("/:id", controller.DeleteUser)
		userRoutes.POST("/:id/activate", controller.ActivateUser)
		userRoutes.POST("/:id/deactivate", controller.DeactivateUser)
	}
	return nil
}
