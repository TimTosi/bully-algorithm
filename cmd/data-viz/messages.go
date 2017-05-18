package main

// -----------------------------------------------------------------------------

// Message is a `struct` used for communication between `bully.Bully`s and the
// `DataServer`.
type Message struct {
	NodeID   int  `json:"nodeID"`
	LeaderID int  `json:"leaderID"`
	State    bool `json:"state"`
}
