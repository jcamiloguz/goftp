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

// func TestHandle(t *testing.T) {

// 	// Test Cases
// 	tables := []struct {
// 		msg    []byte
// 		expect client.ACTIONID
// 	}{
// 		{[]byte("register\n"), client.REG},
// 		{[]byte("ERR\n"), client.ERR},
// 		{[]byte("out\n"), client.OUT},
// 		{[]byte("publish\n"), client.PUB},
// 		{[]byte("subscribe\n"), client.SUB},
// 		{[]byte("subscribe channel=1\n"), client.SUB},
// 		{[]byte("unsubscribe\n"), client.UNSUB},
// 		{[]byte("\n"), client.ERR},
// 	}

// 	for _, item := range tables {
// 		serverMock, connMock := net.Pipe()
// 		defer serverMock.Close()
// 		defer connMock.Close()

// 		s, err := server.NewServer(&server.Config{
// 			Host:      "localhost",
// 			Port:      "3090",
// 			NChannels: 3,
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		clientTest, _ := client.NewClient(connMock, s.Actions, s.Response)

// 		defer serverMock.Close()
// 		wg := sync.WaitGroup{}
// 		go func() {
// 			serverMock.Write(item.msg)
// 		}()
// 		wg.Add(1)
// 		go func() {
// 			select {
// 			case action := <-s.Actions:
// 				if action.Id != item.expect {
// 					t.Errorf("expected %d, got %d", item.expect, action.Id)
// 				}
// 				wg.Done()
// 			}
// 		}()

// 		go clientTest.Read()
// 		wg.Wait()

// 	}

// }
