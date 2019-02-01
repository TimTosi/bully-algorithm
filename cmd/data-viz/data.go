package main

import (
	"log"
	"net"
	"strconv"
	"strings"
)

// DataServer is a `struct` used for gathering and forwarding `bully.Bully`s
// state.
type DataServer struct {
	*net.UDPConn

	swMap map[int]*SlidingWindow
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
	return &DataServer{UDPConn: sock, swMap: make(map[int]*SlidingWindow)}, nil
}

// Listen reads, unpacks and forwards packet from `ds.UDPConn` to `ds.swMap`.
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
				log.Println(err)
				continue
			}
			leaderID, err := strconv.Atoi(packetData[1])
			if err != nil {
				log.Println(err)
				continue
			}
			// UNPACK POLICY
			if _, ok := ds.swMap[nodeID]; !ok {
				ds.swMap[nodeID] = NewSlidingWindow()
			}
			go ds.swMap[nodeID].Push(NewMessage(nodeID, leaderID, true))
		}
	}
}

// Data returns a slice of `*Message` describing current nodes state.
func (ds *DataServer) Data() []*Message {
	var msgs []*Message

	for k, sw := range ds.swMap {
		msgs = append(msgs, sw.Pull(k))
	}
	return msgs
}
