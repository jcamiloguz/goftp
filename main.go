package main

import (
	"log"
	"os"

	"github.com/jcamiloguz/goftp/internal/server"
	"github.com/joho/godotenv"
)

var (
	nChannels int
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	nChannels = 3

	s, err := server.NewServer(&server.Config{
		Host:      host,
		Port:      port,
		NChannels: 3,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started #%d channels", len(s.Channels))

}
