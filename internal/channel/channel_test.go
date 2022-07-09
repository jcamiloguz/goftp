package channel_test

import (
	"testing"

	"github.com/jcamiloguz/goftp/internal/channel"
)

func TestNewChannel(t *testing.T) {
	// Test cases
	tables := []struct {
		idChannel int
	}{
		{0},
		{1},
	}

	// Run test cases
	for _, table := range tables {
		_, err := channel.NewChannel(table.idChannel)
		if table.idChannel < 1 {
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		}
	}
}

// func TestBroadcast(t *testing.T) {
// 	// Simulate a client connection

// 	// Test cases
// 	tables := []struct {
// 		content []byte
// 	}{
// 		{
// 			content: []byte("Hello World!\n"),
// 		},
// 		{

// 			content: []byte("Hello World!\n"),
// 		},
// 		{

// 			content: nil,
// 		},
// 	}

// 	// Run test cases
// 	for _, item := range tables {

// 		serverMock, clientConn := net.Pipe()

// 		actions := make(chan *client.Action)

// 		testChannel, _ := channel.NewChannel(1)
// 		senderClient, _ := client.NewClient(serverMock, "testName", actions)
// 		subscriberClient, _ := client.NewClient(clientConn, "testName", actions)

// 		defer clientConn.Close()
// 		defer serverMock.Close()
// 		wg := sync.WaitGroup{}
// 		testChannel.Clients[subscriberClient.Id] = subscriberClient
// 		wg.Add(1)

// 		t.Logf("Clients channel : %d", len(testChannel.Clients))
// 		go func(content []byte) {
// 			t.Logf("inside go func")
// 			defer wg.Done()

// 			err := testChannel.Broadcast(senderClient, content, t)
// 			t.Logf("Before ")

// 			if content == nil {
// 				if err == nil {
// 					t.Errorf("expected error, got nil")
// 				}
// 			} else {
// 				if err != nil {
// 					t.Errorf("expected nil, got %v", err)
// 				}
// 			}
// 		}(item.content)
// 		go func() {
// 			t.Logf("inside go func")
// 			defer wg.Done()
// 			action := <-actions
// 			switch action.Id {
// 			case client.PUB:
// 				t.Errorf("expected %d, got %d", client.PUB, action.Id)

// 			}
// 		}()
// 		for {
// 			subscriberClient.Read()
// 		}
// 		wg.Wait()

// 		// for {
// 		// 	n, err := bufio.NewReader(subscriberClient.Connection).ReadBytes('\n')
// 		// 	if err == io.EOF {
// 		// 		break
// 		// 	}

// 		// 	if err == nil || string(n) != string(item.content) {
// 		// 		t.Errorf("expected %s, got %s", string(item.content), string(n))
// 		// 		break
// 		// 	}

// 		// 	break
// 		// }

// 	}

// }
