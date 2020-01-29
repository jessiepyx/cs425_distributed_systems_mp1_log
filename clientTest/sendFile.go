package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func sendFile(host string, index int, channel chan string) {
	content, err := ioutil.ReadFile("fakeLog.out")
	if err != nil {
		//Do something
	}
	lines := strings.Split(string(content), "\n")
	linesTotal := len(lines)
	var message string
	for i := index; i < index+linesTotal/10; i++ {
		message += lines[i]
	}

	port := "5001"
	hostPort := string(host + ":" + port)
	log.Println("Establishing TCP connection with " + hostPort)
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		log.Println("Failed to connect to " + hostPort)
		log.Println(err.Error())
		channel <- "fail"
		return
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(message + "\n")); err != nil {
		log.Println("Failed to send request to " + hostPort)
		log.Println(err.Error())
		channel <- "fail"
		return
	}
	channel <- "ok"
	// Send request
}

func main() {
	// generate fake file
	pattern := "true pattern"

	// check arguments
	// if len(os.Args) != 2 {
	// 	log.Fatal("Invalid arguments")
	// }

	// generate fake log file
	// Create or clear output file
	if _, err := os.Stat("fakeLog.out"); os.IsNotExist(err) {
		if _, err = os.Create("fakeLog.out"); err != nil {
			log.Fatal("Failed to create fake log file")
		}
	} else if err = os.Truncate("fakeLog.out", 0); err != nil {
		log.Fatal("Failed to clear fake log file")
	}

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" + "abcdefghijklmnopqrstuvwxyzåäö")
	length := 8000
	var b strings.Builder
	for j := 0; j < 50; j++ {
		for i := 0; i < length; i++ {
			b.WriteRune(chars[rand.Intn(len(chars))])
		}
		b.WriteRune('\n')
	}
	garbageString := b.String()

	// fmt.Println("garbage string is" + garbageString)
	finalTestString := pattern + garbageString
	if f, err := os.OpenFile("fakeLog.out", os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		log.Println(err.Error())
	} else if _, err = f.WriteString(string(finalTestString)); err != nil {
		log.Println("Failed to write to fake log file")
		log.Println(err.Error())
	}

	var servers []string = []string{
		"fa19-cs425-g32-01.cs.illinois.edu",
		"fa19-cs425-g32-02.cs.illinois.edu",
		"fa19-cs425-g32-03.cs.illinois.edu",
		"fa19-cs425-g32-04.cs.illinois.edu",
		"fa19-cs425-g32-05.cs.illinois.edu",
		"fa19-cs425-g32-06.cs.illinois.edu",
		"fa19-cs425-g32-07.cs.illinois.edu",
		"fa19-cs425-g32-08.cs.illinois.edu",
		"fa19-cs425-g32-09.cs.illinois.edu",
		"fa19-cs425-g32-10.cs.illinois.edu",
	}

	var hosts []string
	// Find localhost IP
	addrs, _ := net.InterfaceAddrs()
	var curIP string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			curIP = ipnet.IP.String()
			log.Println("Localhost IP: " + curIP)
		}
	}
	for _, server := range servers {
		command := exec.Command("/usr/bin/dig", "+short", server)
		host, err := command.Output()
		if err != nil {
			log.Println("DNS failed for server " + server)
			log.Println(err.Error())
		} else if ip := string(bytes.Trim(host, "\n")); ip != curIP {
			hosts = append(hosts, ip)
			log.Println("Adding " + ip + " to host list")
		}
	}

	// Query all remote hosts
	channel := make(chan string)
	segment := 0
	for _, host := range hosts {
		go sendFile(host, segment, channel)

		segment++
	}
	var msg string
	for i := 0; i < 9; i++ {
		msg = <-channel
	}
	if msg == "ok" {
		fmt.Println("sent file")
	} else {
		fmt.Println("faied to send file")
	}
}
