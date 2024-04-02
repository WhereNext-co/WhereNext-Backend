package authController

import (
	"net/http"

	authService "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth/service"
	"github.com/gin-gonic/gin"
)

type AuthControllerInterface interface {
    CreateFirebaseUser(c *gin.Context)
}

type authController struct {
    authService authService.AuthServiceInterface
}

func NewAuthController(authService authService.AuthServiceInterface) *authController {
    return &authController{authService}
}

// CreateFirebaseUser creates a new Firebase user
func (uc *authController) CreateFirebaseUser(c *gin.Context) {
    telNo := c.PostForm("telNo")
    if telNo == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Telephone number is required"})
        return
    }

    // Call the CreateFirebaseUser method in the userService
    user, err := uc.authService.CreateFirebaseUser(telNo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"uid": user.UID, "email": user.Email})
}