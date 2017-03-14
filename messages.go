package bully

// -----------------------------------------------------------------------------

// Message Type.
const (
	ELECTION = iota
	OK
	COORDINATOR
	CLOSE
)

// -----------------------------------------------------------------------------

// Message is a structure used within communication between `Bully`s.
type Message struct {
	PeerID string
	Addr   string
	Type   int
}
