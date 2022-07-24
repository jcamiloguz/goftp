package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	cl "github.com/jcamiloguz/goftp/internal/client"
	s "github.com/jcamiloguz/goftp/internal/server"
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

	s, err := s.NewServer(&s.Config{
		Host:      host,
		Port:      port,
		NChannels: 3,
	})
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port))
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	go s.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		newClient, err := cl.NewClient(conn, s.Requests, s.Response)
		if err != nil {
			log.Println(err)
		}

		go func() {
			err := newClient.Read()
			if err != nil {
				log.Println("Error Reading request: ", err.Error())
			}
		}()

	}
}
