package bully

import (
	"fmt"
	"log"
	"net"
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

// mockSocket is a testing function creating a mock socket.
func mockSocket(addr string) (*net.TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, fmt.Errorf("mockSocket: %v", err)
	}
	tcpListener, err := net.ListenTCP("tcp4", laddr)
	if err != nil {
		return nil, fmt.Errorf("mockSocket: %v", err)
	}
	go func() {
		for {
			_, err := tcpListener.AcceptTCP()
			if err != nil {
				log.Printf("listen: %v", err)
				continue
			}
		}
	}()
	return tcpListener, nil
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

func TestBully_Listen(t *testing.T) {
	testCases := []struct {
		name               string
		mockProto          string
		mockAddr           string
		expectedAssertFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			"regular",
			"tcp4",
			"127.0.0.1:8100",
			assert.Nil,
		},
		{
			"badProto",
			"tcp22",
			"127.0.0.1:8101",
			assert.NotNil,
		},
		{
			"badAddr",
			"tcp6",
			"mockBadAddr:8102",
			assert.NotNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := mockBully("1", "127.0.0.1", "1")
			tc.expectedAssertFunc(t, b.Listen(tc.mockProto, tc.mockAddr))
		})
	}
}

func TestPeer_connect(t *testing.T) {
	testCases := []struct {
		name               string
		mockProto          string
		mockAddr           string
		expectedAssertFunc func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			"regular",
			"tcp4",
			"127.0.0.1:8200",
			assert.Nil,
		},
		{
			"badProto",
			"tcp22",
			"127.0.0.1:8200",
			assert.NotNil,
		},
		{
			"badAddr",
			"tcp6",
			"127.0.0.1:9999",
			assert.NotNil,
		},
	}

	ms, err := mockSocket("127.0.0.1:8200")
	assert.Nil(t, err)
	defer func() { _ = ms.Close() }()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := mockBully("1", "127.0.0.1", "1")
			tc.expectedAssertFunc(t, b.connect(tc.mockProto, tc.mockAddr, "1"))
		})
	}
}

func TestPeer_Connect(t *testing.T) {
	testCases := []struct {
		name           string
		mockProto      string
		mockSocketAddr string
		mockPeers      map[string]string
	}{
		{
			"regular",
			"tcp4",
			"127.0.0.1:8300",
			map[string]string{
				"1": "127.0.0.1:8301",
				"2": "127.0.0.1:8302",
				"3": "127.0.0.1:8303",
			},
		},
		{
			"samePeers",
			"tcp4",
			"127.0.0.1:8310",
			map[string]string{
				"1": "127.0.0.1:8311",
				"2": "127.0.0.1:8311",
				"3": "127.0.0.1:8311",
			},
		},
		{
			"emptyMap",
			"tcp4",
			"127.0.0.1:8320",
			map[string]string{},
		},
		{
			"badProto",
			"tcp22",
			"127.0.0.1:8330",
			map[string]string{
				"1": "127.0.0.1:8331",
				"2": "127.0.0.1:8332",
				"3": "127.0.0.1:8333",
			},
		},
		{
			"badPeer",
			"tcp22",
			"127.0.0.1:8340",
			map[string]string{
				"1": "127.0.0.1:8341",
				"2": "notWorkingAddr",
				"3": "127.0.0.1:8343",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ms, err := mockSocket(tc.mockSocketAddr)
			assert.Nil(t, err)
			defer func() { _ = ms.Close() }()

			b := mockBully("1", "127.0.0.1", "1")
			assert.NotPanics(t, func() { b.Connect(tc.mockProto, tc.mockPeers) })
		})
	}
}

func TestBully_SetCoordinator(t *testing.T) {
	testCases := []struct {
		name                string
		mockID              string
		mockPeerID          string
		expectedCoordinator string
	}{
		{"greater", "A", "B", "B"},
		{"less", "Zawarudo", "A", "Zawarudo"},
		{"equal", "same-id", "same-id", "same-id"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := mockBully(tc.mockID, tc.mockID, "127.0.0.1")
			b.SetCoordinator(tc.mockPeerID)
			assert.Equal(t, tc.expectedCoordinator, b.coordinator)
		})
	}
}
