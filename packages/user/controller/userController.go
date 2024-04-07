package userController

import (
	"net/http"

	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"
	"github.com/gin-gonic/gin"
)

type UserControllerInterface interface {
	CreateUserInfo(c *gin.Context)
	CheckUserName(c *gin.Context)
}	

type UserController struct {
	userService userService.UserServiceInterface
}

func NewUserController(userService userService.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

//CreateUserInfo at database
func (uc *UserController) CreateUserInfo(c *gin.Context) {
    var user struct {
        UserName        string `json:"userName"`
        Email           string `json:"email"`
        Title           string `json:"title"`
        Name            string `json:"name"`
        Birthdate       string `json:"birthdate"`
        Region          string `json:"region"`
        TelNo           string `json:"telNo"`
        ProfilePicture  string `json:"profilePicture"`
        Bio             string `json:"bio"`
    }

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    err := uc.userService.CreateUserInfo(user.UserName, user.Email, user.Title, user.Name, user.Birthdate, user.Region, user.TelNo, user.ProfilePicture, user.Bio)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (uc *UserController) CheckUserName(c *gin.Context) {
    var request struct {
        UserName string `json:"userName"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    exists, err := uc.userService.CheckUserName(request.UserName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"exists": exists})
}
