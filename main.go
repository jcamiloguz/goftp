package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	cl "github.com/jcamiloguz/goftp/internal/client"
	s "github.com/jcamiloguz/goftp/internal/server"
	"github.com/jcamiloguz/goftp/internal/webscoket"
	"github.com/joho/godotenv"
)

var (
	nChannels = flag.Int("nchannels", 3, "Number of channels")
)

func getEnv() string {
	return os.Getenv("APP_ENV")
}

func main() {
	env := getEnv()

	if env == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	flag.Parse()

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	s, err := s.NewServer(&s.Config{
		Host:      host,
		Port:      port,
		NChannels: *nChannels,
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
	go webscoket.Start(s.Outbound, s.Inbound)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		newClient, err := cl.NewClient(conn, s.Requests, s.Response)
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			err := newClient.Read()
			if err != nil {
				log.Println("Error Reading request: ", err.Error())
			}
		}()

	}
}
