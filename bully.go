package bully

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Bully is a `struct` representing a single node used by the `Bully Algorithm`.
//
// NOTE: More details about the `Bully algorithm` can be found here
// https://en.wikipedia.org/wiki/Bully_algorithm .
type Bully struct {
	*net.TCPListener

	ID           string
	addr         string
	coordinator  string
	peers        Peers
	mu           *sync.RWMutex
	receiveChan  chan Message
	electionChan chan Message
}

// NewBully returns a new `Bully` or an `error`.
//
// NOTE: All connections to `Peer`s are established during this function.
//
// NOTE: The `proto` value can be one of this list: `tcp`, `tcp4`, `tcp6`.
func NewBully(ID, addr, proto string, peers map[string]string) (*Bully, error) {
	b := &Bully{
		ID:           ID,
		addr:         addr,
		coordinator:  ID,
		peers:        NewPeerMap(),
		mu:           &sync.RWMutex{},
		electionChan: make(chan Message, 1),
		receiveChan:  make(chan Message),
	}

	if err := b.Listen(proto, addr); err != nil {
		return nil, fmt.Errorf("NewBully: %v", err)
	}

	b.Connect(proto, peers)
	return b, nil
}

// receive is a helper function handling communication between `Peer`s
// and `b`. It creates a `gob.Decoder` and a from a `io.ReadCloser`. Each
// `Message` received that is not of type `CLOSE` is pushed to `b.receiveChan`.
//
// NOTE: this function is an infinite loop.
func (b *Bully) receive(rwc io.ReadCloser) {
	var msg Message
	dec := gob.NewDecoder(rwc)

	for {
		if err := dec.Decode(&msg); err == io.EOF || msg.Type == CLOSE {
			_ = rwc.Close()
			if msg.PeerID == b.Coordinator() {
				b.peers.Delete(msg.PeerID)
				b.SetCoordinator(b.ID)
				b.Elect()
			}
			break
		} else if msg.Type == OK {
			select {
			case b.electionChan <- msg:
				continue
			case <-time.After(10 * time.Millisecond):
				continue
			}
		} else {
			b.receiveChan <- msg
		}
	}
}

// listen is a helper function that spawns goroutines handling new `Peers`
// connections to `b`'s socket.
//
// NOTE: this function is an infinite loop.
func (b *Bully) listen() {
	for {
		conn, err := b.AcceptTCP()
		if err != nil {
			log.Printf("listen: %v", err)
			continue
		}
		go b.receive(conn)
	}
}

// Listen makes `b` listens on the address `addr` provided using the protocol
// `proto` and returns an `error` if something occurs.
func (b *Bully) Listen(proto, addr string) error {
	laddr, err := net.ResolveTCPAddr(proto, addr)
	if err != nil {
		return fmt.Errorf("Listen: %v", err)
	}
	b.TCPListener, err = net.ListenTCP(proto, laddr)
	if err != nil {
		return fmt.Errorf("Listen: %v", err)
	}
	go b.listen()
	return nil
}

// connect is a helper function that resolves the tcp address `addr` and try
// to establish a tcp connection using the protocol `proto`. The established
// connection is set to `b.peers[ID]` or the function returns an `error`
// if something occurs.
func (b *Bully) connect(proto, addr, ID string) error {
	raddr, err := net.ResolveTCPAddr(proto, addr)
	if err != nil {
		return fmt.Errorf("connect: %v", err)
	}
	sock, err := net.DialTCP(proto, nil, raddr)
	if err != nil {
		return fmt.Errorf("connect: %v", err)
	}

	b.peers.Add(ID, addr, sock)
	return nil
}

// Connect performs a connection to the remote `Peer`s.
func (b *Bully) Connect(proto string, peers map[string]string) {
	for ID, addr := range peers {
		if b.ID == ID {
			continue
		}
		if err := b.connect(proto, addr, ID); err != nil {
			continue
		}
	}
}

// Send sends a `bully.Message` of type `what` to `b.peer[to]` at the address
// `addr`. If no connection is reachable at `addr` or if `b.peer[to]` does not
// exist, the function retries five times and returns an `error` if it does not
// succeed.
func (b *Bully) Send(to, addr string, what int) error {
	maxRetries := 5

	if !b.peers.Find(to) {
		_ = b.connect("tcp4", addr, to)
	}

	for attempts := 0; ; attempts++ {
		err := b.peers.Write(to, &Message{PeerID: b.ID, Addr: b.addr, Type: what})
		if err == nil {
			break
		}
		if attempts > maxRetries && err != nil {
			return fmt.Errorf("Send: %v", err)
		}
		_ = b.connect("tcp4", addr, to)
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

// SetCoordinator sets `ID` as the new `b.coordinator` if `ID` is greater than
// `b.coordinator` or equal to `b.ID`.
//
// NOTE: This function is thread-safe.
func (b *Bully) SetCoordinator(ID string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if ID > b.coordinator || ID == b.ID {
		b.coordinator = ID
	}
}

// Coordinator returns `b.coordinator`.
//
// NOTE: This function is thread-safe.
func (b *Bully) Coordinator() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.coordinator
}

// Elect handles the leader election mechanism of the `Bully algorithm`.
func (b *Bully) Elect() {
	for _, rBully := range b.peers.PeerData() {
		if rBully.ID > b.ID {
			_ = b.Send(rBully.ID, rBully.Addr, ELECTION)
		}
	}

	select {
	case <-b.electionChan:
		return
	case <-time.After(time.Second):
		b.SetCoordinator(b.ID)
		for _, rBully := range b.peers.PeerData() {
			if rBully.ID < b.ID {
				_ = b.Send(rBully.ID, rBully.Addr, COORDINATOR)
			}
		}
		return
	}
}

// Run launches the two main goroutine. The first one is tied to the
// `Bully algorithm` while the other one is the execution of `workFunc`.
func (b *Bully) Run(workFunc func()) chan error {
	errChan := make(chan error)

	go func() {
		b.Elect()

		for msg := range b.receiveChan {
			if msg.Type == ELECTION {
				if msg.PeerID < b.ID {
					_ = b.Send(msg.PeerID, msg.Addr, OK)
					_ = b.Send(msg.PeerID, msg.Addr, COORDINATOR)
				} else {
					_ = b.Send(msg.PeerID, msg.Addr, OK)
					b.Elect()
				}
			} else if msg.Type == COORDINATOR {
				b.SetCoordinator(msg.PeerID)
			}
		}
	}()

	go workFunc()

	return errChan
}
