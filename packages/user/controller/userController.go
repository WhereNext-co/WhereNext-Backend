package userController

import (
	"log"
	"net/http"

	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"
	"github.com/gin-gonic/gin"
)

type UserControllerInterface interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	FindUser(c *gin.Context)
}	

type UserController struct {
	userService userService.UserServiceInterface
}

func NewUserController(userService userService.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create a user
}

func (uc *UserController) FindUser(c *gin.Context) {
	log.Println(c.Query("email"))
	user := uc.userService.FindUser(c.Query("email"))
	if user.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user})
}