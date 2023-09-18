package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateToken(username string, role string) (string, error) {
	//creating jwt auth token
	secretKey := []byte(os.Getenv("jwtSuperKey"))
	token := jwt.New(jwt.SigningMethodHS256) //newtokencreation

	//setting payload for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() //token expiry setting

	//signing the token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func UserAuth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "No token was awailable"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else {
		tokenStr := strings.Replace(tokenString, "Bearer ", "", -1)
		//parsing the token
		signMethod := jwt.SigningMethodHS256
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if token.Method != signMethod {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("jwtSuperKey")), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token.Valid {
			// Token is valid, proceed with further processing
			claims, ok := token.Claims.(jwt.MapClaims)
			fmt.Println(claims["username"])
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid token"})
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			if claims["role"] != "user" {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "not a user"})
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Set("username", claims["username"].(string))
			c.Next()
			return
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			// Check the error type
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Token has expired"})
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Token validation error"})
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Token validation error"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}

}
