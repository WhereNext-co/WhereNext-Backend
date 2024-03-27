package auth

// import (
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/WhereNext-co/WhereNext-Backend.git/database"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt"
// 	"golang.org/x/crypto/bcrypt"
// )

// var hmacSampleSecret []byte

// // Binding from JSON
// type RegisterBody struct {
// 	Email    string `json:"email" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// 	TelNo    string `json:"tel_no" binding:"required"`
// }

// func Register(c *gin.Context) {
// 	var json RegisterBody
// 	if err := c.ShouldBindJSON(&json); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	// Check User Exists
// 	var userExist database.UserAuth
// 	database.Db.Where("email = ?", json.Email).First(&userExist)
// 	if userExist.ID > 0 {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Exists"})
// 		return
// 	}
// 	// Create User
// 	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
// 	user := database.UserAuth{Email: json.Email, Password: string(encryptedPassword),
// 		TelNo: json.TelNo}
// 	database.Db.Create(&user)
// 	if user.ID > 0 {
// 		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Create Success", "userId": user.ID})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Create Failed"})
// 	}
// }

// // Binding from JSON
// type LoginBody struct {
// 	Email    string `json:"email" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }

// func Login(c *gin.Context) {
// 	var json LoginBody
// 	if err := c.ShouldBindJSON(&json); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	// Check User Exists
// 	var userExist database.UserAuth
// 	database.Db.Where("email = ?", json.Email).First(&userExist)
// 	if userExist.ID == 0 {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Does Not Exists"})
// 		return
// 	}
// 	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
// 	if err == nil {
// 		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 			"userId": userExist.ID,
// 			"exp":    time.Now().Add(time.Minute * 1).Unix(),
// 		})
// 		// Sign and get the complete encoded token as a string using the secret
// 		tokenString, err := token.SignedString(hmacSampleSecret)
// 		fmt.Println(tokenString, err)

// 		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Login Success", "token": tokenString})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Login Failed"})
// 	}
// }
