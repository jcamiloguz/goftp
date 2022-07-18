package server_test

import (
	"fmt"
	"testing"

	"github.com/jcamiloguz/goftp/internal/server"
)

func TestNewServer(t *testing.T) {
	// Test Cases
	tables := []struct {
		host      string
		port      string
		nChannels int
	}{
		{"localhost", "", 0},
		{"localhost", "8080", 2},
		{"localhost", "8080", 0},
		{"localhost", "8080", 100},
		{"", "8080", 100},
	}
	for _, item := range tables {
		fmt.Printf("Testing with host: %s, port: %s, nChannels: %d\n", item.host, item.port, item.nChannels)
		server, err := server.NewServer(&server.Config{
			Host:      item.host,
			Port:      item.port,
			NChannels: item.nChannels})

		if item.host == "" || item.port == "" || item.nChannels == 0 {
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
			if len(server.Channels) != item.nChannels {
				t.Errorf("expected %d channels, got %d", item.nChannels, len(server.Channels))
			}
			if len(server.Clients) != 0 {
				t.Errorf("expected 0 clients, got %d", len(server.Clients))
			}
			if server.Requests == nil {
				t.Errorf("expected nil, got %v", server.Requests)
			}
		}
	}
}
