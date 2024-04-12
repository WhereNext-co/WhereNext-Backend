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

	scheduleController "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/controller"
	scheduleRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/repo"
	scheduleService "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/service"
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
	scheduleRepo := scheduleRepo.NewScheduleRepo(dbConn)
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
	// Initialize the Schedule services and controllers
	scheduleService := scheduleService.NewScheduleService(scheduleRepo)
	scheduleController := scheduleController.NewScheduleController(scheduleService)
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
	r.GET("/users/profile", userController.FindUserByUid)
	r.PUT("/users/profile", userController.UpdateUserInfo)
	// Friend routes
	r.GET("/users/friends", userController.FriendList)
	r.GET("/users/friends/isfriend", userController.IsFriend)
	r.DELETE("/users/friends", userController.RemoveFriend)
	// Friend request routes
	r.POST("/users/friendrequest", userController.CreateFriendRequest)
	r.PUT("/users/friendrequest", userController.AcceptFriendRequest)
	r.DELETE("/users/friendrequest/decline", userController.DeclineFriendRequest)
	r.DELETE("/users/friendrequest/cancel", userController.CancelFriendRequest)
	r.GET("/users/friendrequest", userController.RequestsReceived)

	// Schedule routes
	r.POST("/schedules/create-schedule", scheduleController.CreateSchedule)
	r.DELETE("/schedules/delete-schedule", scheduleController.DeleteSchedule)
	r.PUT("/schedules/edit-schedule", scheduleController.EditSchedule)
	r.GET("/schedules/get-allschedule", scheduleController.GetSchedule)
	r.GET("/schedules/get-Diary", scheduleController.GetDiary)
	// Invitation routes
	r.PATCH("/schedules/accept-invitation", scheduleController.AcceptInvitation)
	r.PATCH("/schedules/reject-invitation", scheduleController.RejectInvitation)
	port := os.Getenv("PORT")
	r.Run(":" + port)
}

// All New Endpoints
// User Profile Page: GET /users/profile, PUT /users/profile
// Friends Page: GET /users/friends, GET /users/friends/{friend}
// Delete Friend: DELETE /users/friends/{friend}
// Add Friend: POST /users/friendrequest/create/{friend}, PUT /users/friendrequest/accept/{friend}
// Add Friend: DELETE /users/friendrequest/decline/{friend}, DELETE /users/friendrequest/cancel/{friend}
// Friend Request Page: GET /users/friendrequest
