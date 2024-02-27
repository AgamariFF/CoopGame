package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	logInfo := log.New(file, "INFO\t", log.Ldate|log.Ltime)
	logErr := log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	fmt.Println("For exit, enter q")
	fmt.Println("Are you a server or a client? (s/c)")
	var userMes string
	fmt.Scan(&userMes)
	switch userMes {
	case "c":
		ClientGame(logInfo, logErr)
	case "s":
		ServerGame(logInfo, logErr)
	case "q":
		logInfo.Fatalln("The programm is closed without select client/server")
	}
}

func handle(con net.Conn, logInfo, logErr *log.Logger) {
	logInfo.Println("Startin new handle gorutin")
	reader := bufio.NewReader(con)
	message := make([]byte, 1028)
	for {
		len, err := reader.Read(message)
		if err != nil && len > 0 {
			logErr.Fatalln(err)
		} else if len > 0 {
			logInfo.Printf("A new message has been received from %q content %q", con.LocalAddr(), string(message[:len]))
		}
		if len > 0 {
			fmt.Println(string(message[:len]))
			len = 0
		}
	}
}

func ServerGame(logInfo, logErr *log.Logger) {
	logInfo.SetPrefix("SERVER_INFO\t")
	logErr.SetPrefix("SERVER_ERROR\t")
	stream, err := net.Listen("tcp", ":8080")
	defer stream.Close()
	if err != nil {
		logErr.Fatal(err)
	} else {
		logInfo.Println("stream has started")
	}
	for {

		go func() {
			reader := bufio.NewReader(os.Stdin)
			userMes := make([]byte, 1)
			for {
				_, err = reader.Read(userMes)
				if err != nil {
					logErr.Println(err)
				}
				if userMes == "q" {
					logInfo.Fatalln("Server closed")
					return
			
				}
			}
		}()

		con, err := stream.Accept()
		if err != nil {
			logErr.Fatal(err)
		} else {
			logInfo.Printf("Accept new con - %q", con.LocalAddr())
		}

		go handle(con, logInfo, logErr)
	}
}

func ClientGame(logInfo, logErr *log.Logger) {
	logInfo.SetPrefix("CLIENT_INFO\t")
	logErr.SetPrefix("CLIENT_ERROR\t")
	conn, err := net.Dial("tcp", ":8080")
	defer conn.Close()
	defer logInfo.Println("Dial has been ending")
	if err != nil {
		logErr.Fatal(err)
	} else {
		logInfo.Println("Dial has started")
	}
	reader := bufio.NewReader(os.Stdin)
	message := make([]byte, 1028)
	for {
		len, err := reader.Read(message)
		if err != nil && len > 0 {
			logErr.Fatalln(err)
		} else if len > 0 {
			logInfo.Printf("A new message has been received from os.stdin, content %q", string(message[:len]))
		}
		SendMessage(message, len, conn, logInfo, logErr)
	}
}

func SendMessage(message []byte, len int, conn net.Conn, logInfo, logErr *log.Logger) {
	len, err := conn.Write(message[:len])
	if err != nil && len > 0 {
		logErr.Fatalln(err)
	} else if len > 0 {
		logInfo.Printf("A new message has been shipped to %q, content %q", conn.LocalAddr(), string(message[:len]))
	}
}
