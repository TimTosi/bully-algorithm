package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/viper"
	bully "github.com/timtosi/bully-algorithm"
)

// -----------------------------------------------------------------------------

// makeBully is a helper function which main purpose is code readability.
func makeBully() (*bully.Bully, error) {
	return bully.NewBully(
		os.Args[1],
		viper.GetStringMapString(confPeerAddr)[os.Args[1]],
		"tcp4",
		viper.GetStringMapString(confPeerAddr),
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

// -----------------------------------------------------------------------------

func main() {
	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	b, err := makeBully()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := getConn()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	workFunc := func() {
		for {
			fmt.Fprintf(conn, "%s:%s", b.ID, b.Coordinator())
			fmt.Printf("Bully %s: Coordinator is %s\n", b.ID, b.Coordinator())
			time.Sleep(1 * time.Second)
		}
	}

	if err := <-b.Run(workFunc); err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown.")
}
