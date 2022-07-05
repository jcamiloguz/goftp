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
