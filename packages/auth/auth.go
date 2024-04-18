package auth

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Client struct {
	FirebaseAuthClient *auth.Client
}

var (
	App        *firebase.App
	AuthClient *Client
)

func InitializeFirebase() (*Client, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("packages/auth/wherenext-24624-firebase-adminsdk-accountkey.json")
	//opt := option.WithCredentialsFile("/run/secrets/firebase_credentials")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	firebaseAuthClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	AuthClient = &Client{FirebaseAuthClient: firebaseAuthClient}

	return AuthClient, nil
}
