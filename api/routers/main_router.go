package routers

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/metsasft/config"
	"github.com/hunderaweke/metsasft/database"
)

func SetupRoutes(r *gin.Engine, c config.Config) error {
	db, err := database.NewMongoDatabase(c)
	if err != nil {
		return err
	}
	ctx := context.Background()
	AddUserRoutes(r, db, ctx)
	AddBlogRoutes(r, db, ctx)
	// AddAuthRoutes(r, db, ctx)
	return nil
}

func Run() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	err = SetupRoutes(r, config)
	if err != nil {
		log.Fatal(err)
	}
	r.Run(":" + config.Server.Port)
}
