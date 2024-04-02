package authService

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	firebaseAuth "firebase.google.com/go/v4/auth"
	auth "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth"
)

type AuthServiceInterface interface {
    CreateFirebaseUser(email string) (*firebaseAuth.UserRecord, error)
}

type authService struct {
    authClient *auth.Client
}

func NewAuthService(authClient *auth.Client) *authService {
    return &authService{authClient}
}

func (us *authService) CreateFirebaseUser(telNo string) (*firebaseAuth.UserRecord, error) {
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
    client := us.authClient

    return client.FirebaseAuthClient.CreateUser(ctx, params)
}