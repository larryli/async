package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

// Current version of ASync
var Version = "0.0.0-dev" // Set at build time in the Makefile

var printVersion = flag.Bool("version", false, "Print version and exit")
var listenAddr = flag.String("listenAddr", "localhost:8282", "Listen address for HTTP server")
var listenNetwork = flag.String("listenNetwork", "tcp", "Listen 'network' (tcp, tcp4, tcp6, unix)")
var listenUmask = flag.Int("listenUmask", 022, "Umask for Unix socket, default: 022")
var cmdPrefix = flag.String("cmdPrefix", "", "Prefix of the command, like /usr/local/bin/php")
var cmdsChanSize = flag.Int("cmdChanSize", 64, "Size of the commands chan, default: 64")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  %s [OPTIONS]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	version := fmt.Sprintf("async %s", Version)
	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Printf("Starting %s", version)

	// Good housekeeping for Unix sockets: unlink before binding
	if *listenNetwork == "unix" {
		if err := os.Remove(*listenAddr); err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}

	// Change the umask only around net.Listen()
	oldUmask := syscall.Umask(*listenUmask)
	listener, err := net.Listen(*listenNetwork, *listenAddr)
	syscall.Umask(oldUmask)
	if err != nil {
		log.Fatal(err)
	}

	async := NewASync(*cmdPrefix, *cmdsChanSize)

	log.Fatal(http.Serve(listener, async))
}
