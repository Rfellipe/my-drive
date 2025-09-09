package main

import (
	"fmt"
	"os"
	"my-drive/lib"
	"my-drive/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func checkEnvVars() {
	vars := []string{
		"DB_USER",
    "DB_PASS",
    "DB_NAME",
		"JWT_SECRET",
	}

	for i, v := range vars {
		variable := os.Getenv(v)
		if variable == "" {
			panic(fmt.Sprintf("variable %s should be in .env file", vars[i]))
		}
	}
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Check for all needed env vars
	checkEnvVars()

	// Get database connection
	lib.DB = lib.StartDatabaseConnection()

	// Start router
	r := gin.Default()
	routes.Auth(r) // Start authentication routes
	routes.File(r) // Start file routes
	r.Run(":8080") // Listen and Server in 0.0.0.0:8080
}
