package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	file, err := os.OpenFile("infoClient.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	logInfo := log.New(file, "INFO_Client\t", log.Ldate|log.Ltime)
	logErr := log.New(file, "ERROR_Client\t", log.Ldate|log.Ltime|log.Lshortfile)
	conn, err := net.Dial("tcp", ":8080")
	defer conn.Close()
	if err != nil {
		logErr.Fatal(err)
	} else {
		logInfo.Println("Dial has started")
	}
	reader := bufio.NewReader(os.Stdin)
	message := make([]byte,1028)
	for {
		len, err := reader.Read(message)
		if err != nil && len > 0 {
			logErr.Fatalln(err)
		} else if len > 0 {
			logInfo.Printf("A new message has been received from os.stdin, content %q", string(message[:len]))
		}
		len, err = conn.Write(message[:len])
		if err != nil && len > 0 {
			logErr.Fatalln(err)
		} else if len > 0 {
			logInfo.Printf("A new message has been shipped to %q, content %q", conn.LocalAddr(), string(message[:len]))
		}
	}
}