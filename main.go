package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type args struct {
	Host string
	Port int
}

type signal struct{}

func main() {
	a := fetchArgs()
	benchmark(a.Host, a.Port)
}

func fetchArgs() *args {
	a := new(args)
	flag.StringVar(&a.Host, "-h", "localhost", "hostname")
	flag.IntVar(&a.Port, "-p", 1234, "portnum")
	flag.Parse()
	return a
}

func benchmark(host string, port int) {
	endSig := make(chan signal)
	go accept(port, endSig)
	if err := dial(host, port); err != nil {
		return
	}
	<-endSig
}

func accept(port int, endSig chan signal) {
	defer close(endSig)
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Addr[%s] listen error: %s\n", addr, err)
		return
	}
	defer ln.Close()

	before := time.Now()
	conn, err := ln.Accept()
	if err != nil {
		fmt.Printf("Addr[%s] accept error: %s\n", addr, err)
		return
	}
	defer conn.Close()
	after := time.Now()

	diff := after.UnixNano() - before.UnixNano()
	fmt.Printf("Addr[%s] accept time: %d nanosecond.\n", addr, diff)
}

func dial(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Addr[%s] dial error: %s\n", addr, err)
		return err
	}
	defer conn.Close()

	return nil
}
