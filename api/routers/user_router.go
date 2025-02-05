package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/metsasft/api/controllers"
	"github.com/hunderaweke/metsasft/api/middlewares"
	"github.com/hunderaweke/metsasft/config"
	"github.com/hunderaweke/metsasft/internal/domain"
	"github.com/hunderaweke/metsasft/internal/usecases"
	"github.com/sv-tools/mongoifc"
)

func AddUserRoutes(r *gin.Engine, db mongoifc.Database, ctx context.Context) error {
	usecase, created, err := usecases.NewUserUsecase(db, ctx)
	if err != nil {
		return err
	}
	if created {
		err = createAdminUser(usecase)
		if err != nil {
			return err
		}
	}
	controller := controllers.NewUserController(usecase)
	r.POST("/login", controller.Login)
	r.POST("/refresh", controller.RefreshToken)
	r.POST("/forgot-password", controller.ForgetPassword)
	r.POST("/reset-password", controller.ResetPassword)
	userRoutes := r.Group("/users")
	userRoutes.GET("/", controller.GetUsers)
	userRoutes.GET("/:id", controller.GetUser)
	userRoutes.Use(middlewares.AuthenticationMiddleware())
	{
		userRoutes.PUT("/:id", controller.UpdateUser)
		userRoutes.POST("/", controller.CreateUser, middlewares.IsAdmin())
		userRoutes.DELETE("/:id", controller.DeleteUser, middlewares.IsAdmin())
		userRoutes.POST("/:id/activate", controller.ActivateUser, middlewares.IsAdmin())
		userRoutes.POST("/:id/deactivate", controller.DeactivateUser, middlewares.IsAdmin())
		userRoutes.POST("/:id/promote", controller.PromoteUser, middlewares.IsAdmin())
		userRoutes.POST("/:id/demote", controller.DemoteUser, middlewares.IsAdmin())
	}
	return nil
}

func createAdminUser(usecase domain.UserUsecase) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}
	adminUser := domain.User{
		Email:            config.Admin.Email,
		Password:         config.Admin.Password,
		TelegramUsername: config.Admin.TelegramUsername,
		IsAdmin:          true,
		IsActive:         true,
	}
	_, err = usecase.CreateUser(adminUser)
	if err != nil {
		return err
	}
	return nil
}
