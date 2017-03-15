package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	bully "github.com/timtosi/bully-algorithm"
)

// -----------------------------------------------------------------------------

func main() {
	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	b, err := bully.NewBully(
		os.Args[1],
		viper.GetStringMapString(confPeerAddr)[os.Args[1]],
		"tcp4",
		viper.GetStringMapString(confPeerAddr),
	)
	if err != nil {
		log.Fatal(err)
	}

	workFunc := func() {
		for {
			fmt.Printf("Bully %s: Coordinator is %s\n", b.ID, b.Coordinator())
			time.Sleep(1 * time.Second)
		}
	}

	if err := <-b.Run(workFunc); err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown.")
}
