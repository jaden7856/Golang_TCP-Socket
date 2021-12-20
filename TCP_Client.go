/* Client n : 1 Server (Read)*/

package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	Conn_Type string = "tcp"
	Conn_Port int    = 8080
)

func main() {
	conn, err := net.Dial(Conn_Type, ":"+strconv.Itoa(Conn_Port))

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go func(c net.Conn) {
		data := make([]byte, 4096)
		for {
			n, err := c.Read(data)
			if err != nil {
				fmt.Println("Fail to read : ", err)
				return
			}
			fmt.Println(string(data[:n]))
		}
	}(conn)

	for {
		var s string
		fmt.Print("입력하세요 : ")
		fmt.Scanln(&s)

		_, err := conn.Write([]byte(s))
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Duration(3) * time.Second)
	}
}
