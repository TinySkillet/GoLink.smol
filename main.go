package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT not found in the environment!")
	}
	server := NewGoLinkServer(":" + portStr)
	server.Run()
}
