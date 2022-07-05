package main

import (
	"flag"
	"log"
	"os"

	"github.com/jcamiloguz/goftp/internal/server"
	"github.com/joho/godotenv"
)

var (
	nChannels = flag.Int("nchannels", 3, "Number of channels")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	s, err := server.NewServer(&server.Config{
		Host:      host,
		Port:      port,
		NChannels: 3,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
	log.Printf("Server started #%d channels", len(s.Channels))

}
