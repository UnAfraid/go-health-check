package main

import (
	"net"
	"fmt"
	"os"
	"time"
	"strconv"
	flag "github.com/spf13/pflag"
)

var timeoutInSeconds int
var version string
var printVersion bool

func init() {
	flag.IntVarP(&timeoutInSeconds, "timeout", "t", 5, "time out in seconds")
	flag.BoolVarP(&printVersion, "version", "v", false, "print the version")
}

func main() {
	flag.Usage = func() {
		fmt.Println("health-check <tcp/udp> <ip> <port>")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if printVersion {
		fmt.Printf("health-check version: %s", version)
		fmt.Println()
		return
	}

	args := flag.Args()
	if len(args) < 3 {
		fmt.Println("health-check <tcp/udp> <ip> <port>")
		os.Exit(1)
	}

	protocol := args[0]
	hostname := args[1]
	port, err := strconv.ParseUint(args[2], 10, 16)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = net.LookupHost(hostname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	timeout := time.Duration(timeoutInSeconds) * time.Second
	conn, err := net.DialTimeout(protocol, fmt.Sprintf("%s:%d", hostname, port), timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't connect to %s://%s:%d: %s\n", protocol, hostname, port, err)
		os.Exit(1)
	}

	err = conn.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't close connection on %s://%s:%d: %s\n", protocol, hostname, port, err)
		os.Exit(2)
	}

	fmt.Printf("%s://%s:%d is alive\n", protocol, hostname, port)
	os.Exit(0)
}
