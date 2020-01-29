package main

import (
	//"bufio"

	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type InMessage struct {
	Query string `json:"query"`
}

type OutMessage struct {
	Total   string `json:"total"`
	Content string `json:"content"`
	Ip      string `json:"ip"`
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
func receiveFakeFile(conn net.Conn) {
	defer conn.Close()
	//receive and write to fake log file
	if _, err := os.Stat("fakeLog.out"); os.IsNotExist(err) {
		if _, err = os.Create("fakeLog.out"); err != nil {
			log.Fatal("Failed to create fake log file")
		}
	} else if err = os.Truncate("fakeLog.out", 0); err != nil {
		log.Fatal("Failed to clear fake log file")
	}
	str, _ := ioutil.ReadAll(conn)
	if f, err := os.OpenFile("fakeLog.out", os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		log.Println(err.Error())
	} else if _, err := f.WriteString(string(str)); err != nil {
		log.Println("Failed to write to fake log file")
		log.Println(err.Error())
	}

}

func main() {
	port := "5001"
	listener, err := net.Listen("tcp", ":"+port)
	checkError(err)
	log.Println("Waiting for connection")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		log.Println("Handling request from client")
		go receiveFakeFile(conn)
	}
}
