package authService

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	firebaseAuth "firebase.google.com/go/v4/auth"
	auth "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth"
	twilio "github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type AuthServiceInterface interface {
    CreateFirebaseUser(email string) (*firebaseAuth.UserRecord, string, error)
    UpdateFirebaseUserPassword(ctx context.Context, telNo string) (string, error)
    SendOTP(telNo string, otp string) error
}

type authService struct {
    authClient *auth.Client
    twilioClient *twilio.RestClient
}

func NewAuthService(authClient *auth.Client, twilioClient *twilio.RestClient) *authService {
    return &authService{authClient, twilioClient}
}

func (us *authService) CreateFirebaseUser(telNo string) (*firebaseAuth.UserRecord, string, error) {
    // Concatenate @wherenext.com to the telephone number
    email := telNo + "@wherenext.com"

    // Generate a random 6-digit password
    randNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
    password := fmt.Sprintf("%06d", randNumber.Intn(1000000))

    params := (&firebaseAuth.UserToCreate{}).
        Email(email).
        Password(password)

    ctx := context.Background()

    // Use the authClient that you've already initialized
    user, err := us.authClient.FirebaseAuthClient.CreateUser(ctx, params)
    if err != nil {
        return nil, "", err
    }

    return user, password, nil
}

// UpdateFirebaseUserPassword updates the password for a Firebase user with the given phone number.
func (us *authService) UpdateFirebaseUserPassword(ctx context.Context, telNo string) (string, error) {
    // Concatenate @wherenext.com to the telephone number
    email := telNo + "@wherenext.com"

    // Generate a new random password
    randNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
    newPassword := fmt.Sprintf("%06d", randNumber.Intn(1000000))

    // Get the user by email
    user, err := us.authClient.FirebaseAuthClient.GetUserByEmail(ctx, email)
    if err != nil {
        return "", fmt.Errorf("error getting user: %v", err)
    }

    // Check if the user was found
    if user == nil {
        return "", fmt.Errorf("no user found with email: %v", email)
    }

    // Update the user's password
    params := (&firebaseAuth.UserToUpdate{}).Password(newPassword)
    _, err = us.authClient.FirebaseAuthClient.UpdateUser(ctx, user.UID, params)
    if err != nil {
        return "", fmt.Errorf("error updating user password: %v", err)
    }

    // Return the new password
    return newPassword, nil
}

func (us *authService) SendOTP(telNo string, otp string) error {
    params := &twilioApi.CreateMessageParams{}
    params.SetTo(telNo)
    params.SetFrom("+12052361785") // Replace with your Twilio number
    params.SetBody("Your OTP is " + otp)

    _, err := us.twilioClient.Api.CreateMessage(params)
    if err != nil {
        return fmt.Errorf("error sending sms message: %v", err)
    }

    return nil
}