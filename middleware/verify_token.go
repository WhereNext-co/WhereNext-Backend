package middleware

import (
	"net/http"
	"strings"

	auth "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth"

	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract the token from the 'Authorization' header
        header := c.GetHeader("Authorization")
        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        token := strings.TrimPrefix(header, "Bearer ")

        // Verify the token using the Firebase Admin SDK
        client := auth.AuthClient.FirebaseAuthClient
        tokenInfo, err := client.VerifyIDToken(c, token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Add the UID to the request context
        c.Set("uid", tokenInfo.UID)

        // Call the next handler
        c.Next()
    }
}