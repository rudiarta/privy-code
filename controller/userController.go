package controller

import (
	"fmt"

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
	userData.Password = HashPassword(c.PostForm("password"))
	db.Create(&userData)

	c.JSON(200, gin.H{
		"message": userData,
	})
}

// HashPassword is function to hash password
func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword)
}

// CheckHashPassword is function to compare hashed password to the original
func CheckHashPassword(hashedPassword, comparePassword string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(comparePassword))
	return result == nil
}
