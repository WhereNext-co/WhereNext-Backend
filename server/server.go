package server

import (
	"log"
	"os"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	auth "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth"
	authController "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth/controller"
	authService "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth/service"
	userController "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/controller"
	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
)

func InitServer() {
	 // Load .env file
	 err := godotenv.Load()
	 if err != nil {
		 log.Fatal("Fail to .env file")
	 }
	dbConn := database.InitDB()
	userRepo := userRepo.NewUserRepo(dbConn)
	 // Initialize Firebase
	 authClient, err := auth.InitializeFirebase()
	 if err != nil {
		 log.Fatalf("error initializing Firebase: %v", err)
	 }

	// Initialize Twilio
	twilioAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioAccountSid,
		Password: twilioAuthToken,
	})

	// Initialize the Authentication services and controllers
	authService := authService.NewAuthService(authClient, twilioClient)
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

	// Auth routes
	r.POST("/auth/createFirebaseUser", authController.CreateFirebaseUser)
	r.POST("/auth/updateFirebaseUserPassword", authController.UpdateFirebaseUserPassword)

	// User routes
	r.POST("/users/create-info", userController.CreateUserInfo)
	r.POST("/users/check-username", userController.CheckUserName)
	port := os.Getenv("PORT")
	r.Run(":"+port)
}

