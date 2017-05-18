package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------

// DataServer is a `struct` used for gathering and forwarding `bully.Bully`s
// state.
type DataServer struct {
	*net.UDPConn

	chanMap map[int]chan Message
}

// NewDataServer returns a new `DataServer` or an `error` if something occurs.
func NewDataServer(addr string) (*DataServer, error) {
	laddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	sock, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}
	return &DataServer{UDPConn: sock, chanMap: make(map[int]chan Message)}, nil
}

// -----------------------------------------------------------------------------

// Listen reads, unpacks and forwards packet from `ds.UDPConn` to `ds.chanMap`.
func (ds *DataServer) Listen() {
	packet := make([]byte, 8)

	for {
		byteRead, _, err := ds.ReadFromUDP(packet)
		if err != nil {
			log.Println(err)
		} else {
			// UNPACK POLICY
			packetData := strings.Split(string(packet[0:byteRead]), ":")
			if len(packetData) < 2 {
				continue
			}
			nodeID, err := strconv.Atoi(packetData[0])
			if err != nil {
				continue
			}
			leaderID, err := strconv.Atoi(packetData[1])
			if err != nil {
				continue
			}
			// UNPACK POLICY
			if _, ok := ds.chanMap[nodeID]; !ok {
				ds.chanMap[nodeID] = make(chan Message)
			}
			select {
			case ds.chanMap[nodeID] <- Message{NodeID: nodeID, LeaderID: leaderID, State: true}:
				continue
			case <-time.After(50 * time.Millisecond):
				continue
			}
		}
	}
}

// TODO
func (ds *DataServer) Data() []Message {
	var msgs []Message

	for k, c := range ds.chanMap {
		select {
		case msg := <-c:
			msgs = append(msgs, msg)
			break
		case <-time.After(1 * time.Second):
			msg := Message{NodeID: k, LeaderID: -1, State: false}
			msgs = append(msgs, msg)
			break
		}
	}

	return msgs
}
