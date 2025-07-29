package main

import (
	"fmt"
	"log"
	"net"
	"strings"
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

		if value.typ != "array" {
			log.Println("Invalid request, expected array length > 0")
			continue
		}

		if len(value.array) == 0 {
			log.Println("Received an empty array!")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		fmt.Println("Result: ", result)
		writer.Write(result)
	}
}
