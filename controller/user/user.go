package user

import (
	"net/http"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"

	"github.com/gin-gonic/gin"
)

func ReadAll(c *gin.Context) {
	var users []database.User
	database.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "users": users})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user database.User
	database.Db.First(&user, userId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user})
}
