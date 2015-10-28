package main

import (
	"fmt"
	"net"
	"reflect"
	"time"
)

type signal struct{}

func main() {
	benchmark("localhost", 1234)
	time.Sleep(time.Second)
	benchmark("127.0.0.1", 2345)
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

	lnType := reflect.TypeOf(ln)
	fmt.Printf("Addr[%s] listener type: %s\n", addr, lnType.String())

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
