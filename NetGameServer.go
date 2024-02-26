package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	file, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	logInfo := log.New(file, "INFO\t", log.Ldate|log.Ltime)
	logErr := log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	stream, err := net.Listen("tcp", ":8080")
	defer stream.Close()
	if err != nil {
		logErr.Fatal(err)
	} else {
		logInfo.Println("stream has started")
	}
	for {
		con, err := stream.Accept()
		if err != nil {
			logErr.Fatal(err)
		} else {
			logInfo.Printf("Accept new con - %s\n", con.LocalAddr())
		}

		go handle(con, *logInfo, *logErr)
	}
}

func handle(con net.Conn, logInfo, logErr log.Logger) {
	logInfo.Println("Startin new handle gorutin")
	reader := bufio.NewReader(con)
	message := make([]byte, 1028)
	for {
		len, err := reader.Read(message)
		if err != nil && len > 0 {
			logErr.Fatalln(err)
		} else if len > 0 {
			logInfo.Printf("A new message has been received from \"%s\" content \"%s\"", con.LocalAddr(), message)
		}
		fmt.Println(string(message[:len]))
	}
}
