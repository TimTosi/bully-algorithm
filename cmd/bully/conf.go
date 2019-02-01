package main

import (
	"log"
	"sort"

	"github.com/spf13/viper"
)

const (
	confFile     = "bully.conf"
	confDataAddr = "data_server_address"
	confPeerAddr = "peer_address"
)

var (
	confPath = []string{"$HOME/", "./", "/tmp/"}
)

func init() {
	viper.SetDefault(confDataAddr, "127.0.0.1:8081")
	viper.SetDefault(confPeerAddr, map[string]string{
		"0": "0.0.0.0:9990",
		"1": "0.0.0.0:9991",
		"2": "0.0.0.0:9992",
		"3": "0.0.0.0:9993",
		"4": "0.0.0.0:9994",
	})

	viper.SetConfigName(confFile)
	for _, p := range confPath {
		viper.AddConfigPath(p)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("'%s': file not found (using defaults).\n", confFile)
	}
	printConf()
}

func printConf() {
	log.Printf("[info] running with configuration found at %v", viper.ConfigFileUsed())
	keys := viper.AllKeys()
	sort.Strings(keys)
	settings := viper.AllSettings()
	for _, k := range keys {
		log.Printf("%v: %+v", k, settings[k])
	}
}
