package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/spf13/viper"
	"github.com/timtosi/bully"
)

// -----------------------------------------------------------------------------

func init() {
	viper.SetDefault("peer_address", map[string]string{
		"0": "0.0.0.0:9990",
		"1": "0.0.0.0:9991",
		"2": "0.0.0.0:9992",
		"3": "0.0.0.0:9993",
		"4": "0.0.0.0:9994",
	})
	viper.SetConfigName("bully.conf")
	for _, p := range []string{"$HOME/", "./", "/tmp/"} {
		viper.AddConfigPath(p)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Println("'bully.conf': file not found (using defaults).")
	}
	printConf()
}

// -----------------------------------------------------------------------------

func printConf() {
	log.Printf("[info] running with configuration found at %v", viper.ConfigFileUsed())
	keys := viper.AllKeys()
	sort.Strings(keys)
	settings := viper.AllSettings()
	for _, k := range keys {
		log.Printf("%v: %+v", k, settings[k])
	}
}

// -----------------------------------------------------------------------------

func main() {

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	b, err := bully.NewBully(
		os.Args[1],
		viper.GetStringMapString("peer_address")[os.Args[1]],
		"tcp4",
		viper.GetStringMapString("peer_address"),
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
