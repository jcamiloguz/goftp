package main

import (
	"log"

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
	nChannels = 3

	s, err := server.NewServer(&server.Config{
		Host:      "localhost",
		Port:      "8080",
		NChannels: 3,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started %d channels", len(s.Channels))

}
