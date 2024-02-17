package main

import (
	AuthController "github.com/WhereNext-co/WhereNext-Backend.git/controller/auth"
	UserController "github.com/WhereNext-co/WhereNext-Backend.git/controller/user"
	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	"github.com/WhereNext-co/WhereNext-Backend.git/initializers"
	"github.com/WhereNext-co/WhereNext-Backend.git/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

// Binding from JSON
type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

func main() {
	database.InitDB()
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	authorized := r.Group("/users", middleware.JWTAuthen())
	authorized.GET("/readall", UserController.ReadAll)
	authorized.GET("/profile", UserController.Profile)
	r.Run(":3000")
}
