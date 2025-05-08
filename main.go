package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, relying on system environment variables")
	}

	SetJWTKey(os.Getenv("JWT_SECRET_KEY"))
	
	InitDB()

	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/refresh", refreshHandler)
	http.HandleFunc("/revoke", authMiddleware(revokeHandler))
	http.HandleFunc("/protected", authMiddleware(protectedHandler))

	fmt.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
