package client_test

import (
	"net"
	"testing"

	"github.com/jcamiloguz/goftp/internal/client"
)

func TestNewClient(t *testing.T) {
	// Simulate a client connection
	_, clientConn := net.Pipe()

	actions := make(chan *client.Action)
	responses := make(chan *client.Action)

	// Test Cases
	tables := []struct {
		conn      net.Conn
		actions   chan *client.Action
		responses chan *client.Action
	}{
		{nil, nil, nil},
		{clientConn, nil, nil},
		{clientConn, nil, responses},
		{clientConn, actions, nil},
		{clientConn, actions, responses},
	}

	for _, item := range tables {
		client, err := client.NewClient(item.conn, item.actions, item.responses)

		if item.conn == nil || item.actions == nil || item.responses == nil {
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
			if client.Connection != item.conn {
				t.Errorf("expected %v, got %v", item.conn, client.Connection)
			}
			if client.StartConnTime.IsZero() {
				t.Errorf("expected non-zero time, got zero")
			}
			if !client.EndConnTime.IsZero() {
				t.Errorf("expected zero time, got non-zero")
			}
		}
	}

}
