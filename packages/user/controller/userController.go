package userController

import (
	"log"
	"net/http"

	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"
	"github.com/gin-gonic/gin"
)

type UserControllerInterface interface {
	CreateUser(c *gin.Context)
	FindUser(c *gin.Context)
	CreateUserInfo(c *gin.Context)
}	

type UserController struct {
	userService userService.UserServiceInterface
}

func NewUserController(userService userService.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	// Implement the logic to create a user
}

func (uc *UserController) CreateUserInfo(c *gin.Context) {
    var user struct {
        Title           string `json:"title"`
        Name            string `json:"name"`
        Birthdate       string `json:"birthdate"`
        Region          string `json:"region"`
        TelNo           string `json:"telNo"`
        ProfilePicture  string `json:"profilePicture"`
    }

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    err := uc.userService.CreateUserInfo(user.Title, user.Name, user.Birthdate, user.Region, user.TelNo, user.ProfilePicture)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
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