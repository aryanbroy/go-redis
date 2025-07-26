package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	input := "$5\r\nAhmed\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	b, err := reader.ReadByte()
	if err != nil {
		log.Println("error reading a bytes from the input: ", err.Error())
		return
	}

	if b != '$' {
		log.Println("Invalid type, expecting bulk strings only!")
		os.Exit(1)
	}

	size, err := reader.ReadByte()
	if err != nil {
		log.Println("Error reading the size byte: ", err.Error())

	}

	strSize, err := strconv.ParseInt(string(size), 10, 64)
	if err != nil {
		log.Fatalln("Error parsing size byte to integer: ", err.Error())
	}

	reader.ReadByte()
	reader.ReadByte()

	name := make([]byte, strSize)
	_, err = reader.Read(name)
	if err != nil {
		log.Fatalln("Error reading data...", err.Error())
	}

	fmt.Println(name)

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
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("IOF error reading into buffer: ", err)
				break
			}
			fmt.Println("Error reading from the client: ", err.Error())
			os.Exit(1)
		}

		conn.Write([]byte("+OK\r\n"))
	}
}
