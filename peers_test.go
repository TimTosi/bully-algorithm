package bully

import (
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
