package bully

import (
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeerMap_NewPeerMap(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{"regular"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, NewPeerMap())
		})
	}
}

func TestPeerMap_Add(t *testing.T) {
	testCases := []struct {
		name      string
		mockPeers []*Peer
	}{
		{
			"regular_single",
			[]*Peer{
				&Peer{ID: "single", addr: "127.0.0.1", sock: gob.NewEncoder(nil)},
			},
		},
		{
			"regular_multiple",
			[]*Peer{
				&Peer{ID: "multiple-1", addr: "40.87.127.215", sock: gob.NewEncoder(nil)},
				&Peer{ID: "multiple-2", addr: "84.72.203.27", sock: gob.NewEncoder(nil)},
				&Peer{ID: "multiple-3", addr: "232.65.164.182", sock: gob.NewEncoder(nil)},
			},
		},
		{"none_added", []*Peer{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			res := NewPeerMap()
			for _, mockPeer := range tc.mockPeers {
				res.Add(mockPeer.ID, mockPeer.addr, nil)
			}
			assert.Equal(t, len(tc.mockPeers), len(res.peers))

			expectedPeerList := make([]*Peer, 0)
			for _, expectedPeer := range res.peers {
				expectedPeerList = append(expectedPeerList, expectedPeer)
			}
			assert.ElementsMatch(t, expectedPeerList, tc.mockPeers)
		})
	}
}

// func TestPeerMap_Delete(t *testing.T) {
// 	testCases := []struct {
// 		name string
// 	}{
// 		{"regular"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			assert.NotNil(t, NewPeerMap())
// 		})
// 	}
// }

// func TestPeerMap_Find(t *testing.T) {
// 	testCases := []struct {
// 		name string
// 	}{
// 		{"regular"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			assert.NotNil(t, NewPeerMap())
// 		})
// 	}
// }

// func TestPeerMap_Write(t *testing.T) {
// 	testCases := []struct {
// 		name string
// 	}{
// 		{"regular"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			assert.NotNil(t, NewPeerMap())
// 		})
// 	}
// }

// func TestPeerMap_PeerData(t *testing.T) {
// 	testCases := []struct {
// 		name string
// 	}{
// 		{"regular"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			assert.NotNil(t, NewPeerMap())
// 		})
// 	}
// }

// ADD RACE TESTS
