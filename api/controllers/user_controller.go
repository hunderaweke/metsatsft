package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hunderaweke/metsasft/internal/domain"
	"github.com/hunderaweke/metsasft/pkg"
)

type UserController struct {
	usecase domain.UserUsecase
}

func NewUserController(usecase domain.UserUsecase) *UserController {
	return &UserController{usecase: usecase}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var user domain.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	user.IsActive = false
	user.IsAdmin = false
	user, err = c.usecase.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	var users []domain.User
	users, err := c.usecase.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.usecase.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
func (c *UserController) UpdateUser(ctx *gin.Context) {
	var user domain.User
	id := ctx.Param("id")
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	user.ID = id
	user, err = c.usecase.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.usecase.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (c *UserController) ActivateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.usecase.ActivateUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *UserController) DeactivateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.usecase.DeactivateUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *UserController) Login(ctx *gin.Context) {
	var user domain.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	user, err = c.usecase.Login(user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	refreshToken, accessToken, err := pkg.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"refresh_token": refreshToken, "access_token": accessToken})
}

func (c *UserController) RefreshToken(ctx *gin.Context) {
	refreshToken := struct {
		Token string `json:"refresh_token"`
	}{}
	if err := ctx.ShouldBindJSON(&refreshToken); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	claims, err := pkg.ValidateRefreshToken(refreshToken.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user, err := c.usecase.GetUserByID(claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newRefreshToken, accessToken, err := pkg.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"refresh_token": newRefreshToken, "access_token": accessToken})
}

func (c *UserController) PromoteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.usecase.PromoteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (c *UserController) DemoteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.usecase.DemoteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
