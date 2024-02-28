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
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To16(); ipv4 != nil && ipv4.String()[:3] == "192" {
			fmt.Println("Your IP: ", ipv4)
		}
	}
	stream, err := net.Listen("tcp", ":3705")
	if err == nil {
		defer stream.Close()
		logInfo.Printf("stream %q has started", stream.Addr().String())
	} else {
		logErr.Fatal(err)
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
				if string(userMes) == "q" {
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
	fmt.Println("Enter the IP address of the server")
	var ip string
	fmt.Scan(&ip)
	ipv4 := ip + ":3705"
	logInfo.Printf("The IP address of the server has been received %q. Trying to connect", ipv4)
	conn, err := net.Dial("tcp", ipv4)
	defer logInfo.Println("Dial has been ending")
	if err != nil {
		logErr.Fatal(err)
	} else {
		defer conn.Close()
		logInfo.Printf("Dial has started with %q", ipv4)
	}
	reader := bufio.NewReader(os.Stdin)
	message := make([]byte, 1028)
	for {
		len, err := reader.Read(message)
		if err != nil && len > 0 {
			logErr.Fatalln(err)
		} else if string(message[:len]) == "q\r\n" {
			logInfo.Fatalln("Exit")
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
