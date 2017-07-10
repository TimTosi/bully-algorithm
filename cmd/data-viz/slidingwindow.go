package main

import "time"

// -----------------------------------------------------------------------------

// SlidingWindow is a `struct` representing a node state over time.
type SlidingWindow struct {
	mc       chan *Message
	lastSent *Message
	id       int
	miss     int
}

// NewSlidingWindow returns a new `*SlidingWindow`.
func NewSlidingWindow() *SlidingWindow { return &SlidingWindow{mc: make(chan *Message, 1)} }

// Push sets the `SlidingWindow` with the latest node state if `sw.mc` is
// available.
func (sw *SlidingWindow) Push(m *Message) {
	select {
	case sw.mc <- m:
		return
	case <-time.After(50 * time.Millisecond):
		return
	}
}

// Pull returns the latest node state available or a `*Message` describing
// the node's dead state if `sw.miss` is greater than `2`.
func (sw *SlidingWindow) Pull(id int) (m *Message) {
	select {
	case m = <-sw.mc:
		sw.miss = 0
		break
	case <-time.After(1 * time.Second):
		if sw.miss > 2 {
			m = NewMessage(id, -1, false)
		} else {
			m = sw.lastSent
			sw.miss++
		}
		break
	}
	sw.lastSent = m
	return m
}
