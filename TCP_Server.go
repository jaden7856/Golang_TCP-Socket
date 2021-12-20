/* Client n : 1 Server (Read)*/

package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	ConnType = "tcp"
	ConnPort = 8080
)

func main() {
	listen, err := net.Listen(ConnType, ":"+strconv.Itoa(ConnPort))

	if err != nil {
		panic(err)
	}
	fmt.Println("Start server!")

	defer listen.Close()

	for {
		conn, err := listen.Accept()

		if err != nil {
			fmt.Println("Fail to Accept; err : ", err)
			continue
		}
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn) {
	//defer conn.Close()
	rcvBuf := make([]byte, 4096)

	for {
		reqLen, err := conn.Read(rcvBuf)

		if err != nil {
			if io.EOF == err {
				fmt.Println("Fail to Read : ", err)
				break
			}
		}

		if reqLen > 0 {
			data := rcvBuf[:reqLen]
			fmt.Println(string(data))
			_, err := conn.Write(data[:reqLen])

			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}
