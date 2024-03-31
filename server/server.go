package server

import (
	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	userController "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/controller"
	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	dbConn := database.InitDB()
	userRepo := userRepo.NewUserRepo(dbConn)
	userService := userService.NewUserService(userRepo)
	userController := userController.NewUserController(userService)

	r := gin.Default()
	r.Use(cors.Default())

	// r.POST("/register", AuthController.Register)
	// r.POST("/login", AuthController.Login)
	// authorized := r.Group("/users", middleware.JWTAuthen())
	// authorized.GET("/readall", UserController.ReadAll)
	// authorized.GET("/profile", UserController.Profile)

	r.GET("/users", userController.FindUser)
	r.POST("/users/create", userController.CreateUser)
	r.POST("/users/create-info", userController.CreateUserInfo)
	r.Run(":3000")
}

