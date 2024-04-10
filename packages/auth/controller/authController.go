package controller

import (
	"context"
	"net/http"

	authService "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth/service"
	"github.com/gin-gonic/gin"
)

type AuthControllerInterface interface {
    CreateFirebaseUser(c *gin.Context)
	UpdateFirebaseUserPassword(c *gin.Context)
}

type authController struct {
    authService authService.AuthServiceInterface
}

func NewAuthController(authService authService.AuthServiceInterface) *authController {
    return &authController{authService}
}

// CreateFirebaseUser creates a new Firebase user
func (uc *authController) CreateFirebaseUser(c *gin.Context) {
    type TelNo struct {
        TelNo string `json:"telNo"`
    }
    var json TelNo
    if err := c.BindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, expected JSON"})
        return
    }
telNo := json.TelNo
    if telNo == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Telephone number is required"})
        return
    }

    // Call the CreateFirebaseUser method in the userService
    user, password, err := uc.authService.CreateFirebaseUser(telNo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    // Send OTP via SMS using Twilio
    err = uc.authService.SendOTP(telNo, password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending OTP"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"uid": user.UID, "email": user.Email, "password": password})
}

// UpdateUserPassword updates a user's password and sends an OTP
func (uc *authController) UpdateFirebaseUserPassword(c *gin.Context) {
    type TelNo struct {
        TelNo string `json:"telNo"`
    }
    var json TelNo
    if err := c.BindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, expected JSON"})
        return
    }
telNo := json.TelNo
    if telNo == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Telephone number is required"})
        return
    }

    // Call the UpdateFirebaseUserPassword method in the authService
    newPassword, err := uc.authService.UpdateFirebaseUserPassword(context.Background(), telNo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user password"})
        return
    }

    // Send OTP via SMS using Twilio
    err = uc.authService.SendOTP(telNo, newPassword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending OTP"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"telNo": telNo, "message": "Password updated and OTP sent"})
}