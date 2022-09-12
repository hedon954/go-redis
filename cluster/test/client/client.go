package main

import (
	"fmt"
	"net"
)

func main() {
	write(":6379")
	read(":6380")
}

func write(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("*3\r\n$3\r\nSET\r\n$3\r\nke1\r\n$3\r\nva1\r\n"))
	bs := make([]byte, 1024)
	conn.Read(bs)

	fmt.Println("listen to:", addr, "got", string(bs))
}

func read(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nke1\r\n"))
	bs := make([]byte, 1024)
	conn.Read(bs)

	fmt.Println("listen to:", addr, "got", string(bs))
}
