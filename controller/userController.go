package controller

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rudiarta/privy-code/database"
	"github.com/rudiarta/privy-code/model"
	"golang.org/x/crypto/bcrypt"
)

// AddUserController is function to handle controller
func AddUserController(c *gin.Context) {
	db, err := database.InitDatabase()
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		db.Close()
		if r := recover(); r != nil {
			fmt.Println(r.(error).Error())
		}
	}()

	userData := model.User{}
	userData.Name = c.PostForm("name")
	userData.Email = c.PostForm("email")
	userData.Password = hashPassword(c.PostForm("password"))
	if err := db.Create(&userData).Error; err != nil {
		c.JSON(422, gin.H{
			"message": "Data duplicate, or data can't be inserted",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": userData,
	})

	return
}

// LoginUserController is function to handle controller
func LoginUserController(c *gin.Context) {
	db, err := database.InitDatabase()
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		db.Close()
	}()

	email := c.PostForm("email")
	password := c.PostForm("password")

	userData := model.User{}
	db.Where("email = ?", email).First(&userData)
	userID := strconv.Itoa(userData.ID)
	if validatePassword := checkHashPassword(userData.Password, password); validatePassword {
		c.JSON(200, gin.H{
			"message": "Success Login.",
			"toke":    generateJwtToken(userID),
		})
		return
	}
	c.JSON(422, gin.H{
		"message": "Wrong Password. " + userData.Email,
	})
	return
}

// hashPassword is function to hash password
func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword)
}

// checkHashPassword is function to compare hashed password to the original
func checkHashPassword(hashedPassword, comparePassword string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(comparePassword))
	return result == nil
}

// generateJwtToken is a function to generate JWT Token
// and will use user_id
func generateJwtToken(userID string) string {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = string(userID)
	atClaims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return ""
	}

	return token
}
