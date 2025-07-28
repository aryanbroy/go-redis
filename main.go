package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error establishing a tcp connection: ", err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error receiving requests: ", err)
		return
	}
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println("Error reading inputs from client: ", err)
			return
		}

		fmt.Println("Value: ", value)

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
