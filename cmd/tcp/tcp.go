package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/donnie4w/go-logger/logger"
)

func vineFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("vine", flag.ExitOnError)
	flagSet.String("vined-address", ":4211", "vined server address")
	flagSet.String("portname", "", "serial portname")
	return flagSet
}

func main() {
	flagSet := vineFlagSet()

	flagSet.Parse(os.Args[1:])
	vinedAddress := flagSet.Lookup("vined-address").Value.String()
	portName := flagSet.Lookup("portname").Value.String()
	fmt.Println(vinedAddress, portName)

	conn, err := net.Dial("tcp", vinedAddress)
	if err != nil {
		log.Error("net.Dial:", err)
		os.Exit(1)
	}
	portNameMd5 := md5.Sum([]byte(portName))

	n, err := conn.Write(portNameMd5[:])
	if err != nil || n <= 0 {
		log.Error("conn.write", err)
		os.Exit(1)
	}

	localAddress := conn.LocalAddr().String()

	go func() {
		for {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				log.Debug("conn.Write:", err)
				os.Exit(1)
			}
			if n > 0 {
				log.Debug("read:", string(buff[:n]))
			}
		}
	}()
	for {
		time.Sleep(time.Second)
		_, err := conn.Write([]byte(localAddress + ":" + "hello"))
		if err != nil {
			log.Debug("conn.Write:", err)
			os.Exit(1)
		}
	}

}
