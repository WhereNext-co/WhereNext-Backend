package server

import (
	"log"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	auth "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth"
	authController "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth/controller"
	authService "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth/service"
	userController "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/controller"
	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	dbConn := database.InitDB()
	userRepo := userRepo.NewUserRepo(dbConn)
	 // Initialize Firebase
	 authClient, err := auth.InitializeFirebase()
	 if err != nil {
		 log.Fatalf("error initializing Firebase: %v", err)
	 }
	// Initialize the Authentication services and controllers
	authService := authService.NewAuthService(authClient)
	authController := authController.NewAuthController(authService)
	// Initialize the User services and controllers
	userService := userService.NewUserService(userRepo)
	userController := userController.NewUserController(userService)

	r := gin.Default()
	r.Use(cors.Default())

	// r.POST("/register", AuthController.Register)
	// r.POST("/login", AuthController.Login)
	// authorized := r.Group("/users", middleware.JWTAuthen())
	// authorized.GET("/readall", UserController.ReadAll)
	// authorized.GET("/profile", UserController.Profile)

	r.POST("/auth/createFirebaseUser", authController.CreateFirebaseUser)
	r.GET("/users", userController.FindUser)
	r.POST("/users/create", userController.CreateUser)
	r.POST("/users/create-info", userController.CreateUserInfo)
	r.Run(":3000")
}

