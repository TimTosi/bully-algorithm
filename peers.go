package bully

import (
	"fmt"
	"io"
	"sync"
)

// Peers is an `interface` exposing methods to handle communication with other
// `bully.Bully`s.
//
// NOTE: This project offers a default implementation of the `Peers` interface
// that provides basic functions. This will work for the most simple of use
// cases fo exemples, although I strongly recommend you provide your own, safer
// implementation while doing real work.
type Peers interface {
	Add(ID, addr string, fd io.Writer)
	Delete(ID string)
	Find(ID string) bool
	Write(ID string, msg interface{}) error
	PeerData() []struct {
		ID   string
		Addr string
	}
}

// PeerMap is a `struct` implementing the `Peers` interface and representing
// a container of `bully.Peer`s.
type PeerMap struct {
	mu    *sync.RWMutex
	peers map[string]*Peer
}

// NewPeerMap returns a new `bully.PeerMap`.
func NewPeerMap() *PeerMap {
	return &PeerMap{mu: &sync.RWMutex{}, peers: make(map[string]*Peer)}
}

// Add creates a new `bully.Peer` and adds it to `pm.peers` using `ID` as a key.
//
// NOTE: This function is thread-safe.
func (pm *PeerMap) Add(ID, addr string, fd io.Writer) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.peers[ID] = NewPeer(ID, addr, fd)
}

// Delete erases the `bully.Peer` corresponding to `ID` from `pm.peers`.
//
// NOTE: This function is thread-safe.
func (pm *PeerMap) Delete(ID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.peers, ID)
}

// Find returns `true` if `pm.peers[ID]` exists, `false` otherwise.
//
// NOTE: This function is thread-safe.
func (pm *PeerMap) Find(ID string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	_, ok := pm.peers[ID]
	return ok
}

// Write writes `msg` to `pm.peers[ID]`. It returns `nil` or an `error` if
// something occurs.
//
// NOTE: This function is thread-safe.
func (pm *PeerMap) Write(ID string, msg interface{}) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if p, ok := pm.peers[ID]; !ok {
		return fmt.Errorf("Write: peer %s not found in PeerMap", ID)

	} else if err := p.sock.Encode(msg); err != nil {
		return fmt.Errorf("Write: %v", err)
	}
	return nil
}

// PeerData returns a slice of anonymous structures representing a tupple
// composed of a `Peer.ID` and `Peer.addr`.
//
// NOTE: This function is thread-safe.
func (pm *PeerMap) PeerData() []struct {
	ID   string
	Addr string
} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var IDSlice []struct {
		ID   string
		Addr string
	}
	for _, peer := range pm.peers {
		IDSlice = append(IDSlice, struct {
			ID   string
			Addr string
		}{
			peer.ID,
			peer.addr,
		})
	}
	return IDSlice
}
