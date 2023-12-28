package main

import (
	"log"
	"os"
	"testing"
	"time"
)

// const bufSize = 1024 * 1024

// var lis *bufconn.Listener

// func bufDialer(context.Context, string) (net.Conn, error) {
// 	return lis.Dial()
// }

func TestMain(m *testing.M) {
	go Server()
	log.Println("preparing data environment . . . ")
	time.Sleep(2 * time.Second)
	log.Println("starting Tests . . . ")
	ret := m.Run()

	if ret == 0 {
		log.Println("Ended Tests")
		teardown()
	}
	os.Exit(ret)
}

func teardown() {
	log.Println("Teardown")
}

func setup() {
}
