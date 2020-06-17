package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rudiarta/privy-code/routes"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("File .env error/not found!!")
	}
	app := gin.Default()
	routes.InitRoutes(app)
	app.Run()
}
