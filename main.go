package main

import (
	"favon-api/utils"
	"favon-api/web"
	"log"

	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create redis client
	redisClient := utils.CreateRedisClient()
	defer redisClient.Close()

	// Create Logger
	utils.InitLog(os.Stdout, os.Stderr)

	// Load routes and boot up server
	web.Init(redisClient)
}
