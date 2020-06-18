package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is function to handle middleware
func AuthMiddleware(c *gin.Context) {
	if checkTokenAttach := c.Request.Header["Authorization"]; len(checkTokenAttach) == 0 {
		fmt.Println("Tokken not attached")
		c.Abort()
		c.JSON(401, gin.H{
			"message": "Token must be attached.",
		})

		return
	}

	tokenString := c.Request.Header["Authorization"][0]
	tokenBearerArray := strings.Fields(tokenString)
	tokenBearer := tokenBearerArray[1]
	extractedTokenValue, err := verifyJwtToken(tokenBearer)
	if err != nil {
		c.Abort()
		c.JSON(401, gin.H{
			"message": "Token not valid.",
		})

		return
	}

	tokenValue, ok := extractedTokenValue.Claims.(jwt.MapClaims)
	if ok && extractedTokenValue.Valid {
		userID := fmt.Sprintf("%v", tokenValue["user_id"])
		c.Set("user_id", userID)
		c.Next()

		return
	}

	c.Abort()
	c.JSON(401, gin.H{
		"message": "Token not valid.",
	})

	return
}

// verifyJwtToken is for verify the token and exratct token value
func verifyJwtToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
