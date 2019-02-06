package bully

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockBully is a testing function returning a mock `*bully.Bully`.
func mockBully(ID, addr, coordinator string) *Bully {
	return &Bully{
		ID:           ID,
		addr:         addr,
		coordinator:  ID,
		peers:        NewPeerMap(),
		mu:           &sync.RWMutex{},
		electionChan: make(chan Message, 1),
		receiveChan:  make(chan Message),
	}
}

// -----------------------------------------------------------------------------

func TestBully_NewBully(t *testing.T) {
	testCases := []struct {
		name                    string
		mockID                  string
		mockAddr                string
		mockProto               string
		mockPeers               map[string]string
		expectedAssertBullyFunc func(assert.TestingT, interface{}, ...interface{}) bool
		expectedAssertErrorFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			"regular", "1",
			"127.0.0.1:8000",
			"tcp4",
			nil,
			assert.NotNil,
			assert.Nil,
		},
		{
			"badProto",
			"1",
			"127.0.0.1:8001",
			"tcp22",
			nil,
			assert.Nil,
			assert.NotNil,
		},
		{
			"badAddr",
			"1",
			"errorAddr:8002",
			"tcp4",
			nil,
			assert.Nil,
			assert.NotNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewBully(tc.mockID, tc.mockAddr, tc.mockProto, tc.mockPeers)
			tc.expectedAssertBullyFunc(t, res)
			tc.expectedAssertErrorFunc(t, err)
		})
	}
}
