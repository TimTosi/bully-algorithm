package main

// Message is a `struct` used for communication between `bully.Bully`s and the
// `DataServer`.
type Message struct {
	NodeID   int  `json:"nodeID"`
	LeaderID int  `json:"leaderID"`
	State    bool `json:"state"`
}

// NewMessage returns a new `*Message`.
func NewMessage(nodeID, leaderID int, state bool) *Message {
	return &Message{NodeID: nodeID, LeaderID: leaderID, State: state}
}
