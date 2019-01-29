package bully

import (
	"encoding/gob"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNewPeer(t *testing.T) {
	testCases := []struct {
		name         string
		expectedID   string
		expectedAddr string
		expectedPeer *Peer
	}{
		{
			"regular",
			"test-peer",
			"127.0.0.1",
			&Peer{ID: "test-peer", addr: "127.0.0.1", sock: gob.NewEncoder(nil)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := NewPeer(tc.expectedID, tc.expectedAddr, nil)
			assert.Equal(t, res, tc.expectedPeer)
		})
	}
}
