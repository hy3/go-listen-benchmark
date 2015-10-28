package main

import (
	"fmt"
	"net"
	"reflect"
	"time"
)

type signal struct{}

func main() {
	benchmark("localhost:1234")
	benchmark("127.0.0.1:2345")
}

func benchmark(addr string) {
	endSig := make(chan signal)
	go accept(addr, endSig)
	if err := dial(addr); err != nil {
		fmt.Printf("Addr[%s] dial error: %s\n", addr, err)
		return
	}
	<-endSig
}

func accept(addr string, endSig chan signal) {
	defer close(endSig)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Addr[%s] listen error: %s\n", addr, err)
		return
	}

	lnType := reflect.TypeOf(ln)
	fmt.Printf("Addr[%s] listener type: %s\n", addr, lnType.String())

	before := time.Now()
	conn, err := ln.Accept()
	if err != nil {
		fmt.Printf("Addr[%s] listen error: %s\n", addr, err)
		return
	}
	defer conn.Close()
	after := time.Now()

	diff := after.UnixNano() - before.UnixNano()
	fmt.Printf("Addr[%s] accept time: %d nanosecond.\n", addr, diff)
}

func dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}