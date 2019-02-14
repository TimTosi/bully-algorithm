package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	bully "github.com/timtosi/bully-algorithm"
)

// makeBully is a helper function which main purpose is code readability.
func makeBully() (*bully.Bully, error) {
	peers := viper.GetStringMapString(confPeerAddr)

	addr := strings.Split(peers[os.Args[1]], ":")
	if len(addr) < 2 {
		return nil, fmt.Errorf("makeBully: cannot retrieve local address")
	}

	return bully.NewBully(
		os.Args[1],
		fmt.Sprintf("0.0.0.0:%s", addr[1]),
		"tcp4",
		peers,
	)
}

// getConn is a helper function which main purpose is code readability.
func getConn() (*net.UDPConn, error) {
	raddr, err := net.ResolveUDPAddr("udp", viper.GetString(confDataAddr))
	if err != nil {
		return nil, err
	}
	return net.DialUDP("udp", nil, raddr)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("[ERR] node ID required such as `./bully 1`")
	}

	b, err := makeBully()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := getConn()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	workFunc := func() {
		for {
			_, _ = fmt.Fprintf(conn, "%s:%s", b.ID, b.Coordinator())
			fmt.Printf("Bully %s: Coordinator is %s\n", b.ID, b.Coordinator())
			time.Sleep(1 * time.Second)
		}
	}

	b.Run(workFunc)
}
