/* Client n : 1 Server (Read)*/

package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")

	if err != nil {
		panic(err)
	}
	defer conn.Close()

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
