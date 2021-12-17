/* Client n : 1 Server (Read)*/

package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	CONN_Type = "tcp"
	CONN_Port = 8080
)

func main() {
	conn, err := net.Dial(CONN_Type, ":"+strconv.Itoa(CONN_Port))

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
		conn.Write([]byte(s))
		time.Sleep(time.Duration(3) * time.Second)
	}
}
