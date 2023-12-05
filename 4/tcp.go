package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

func handleConnection(conn net.Conn, Mute *sync.Mutex) {
	defer conn.Close()

	fmt.Println("Клиент подключён", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	set := &Set{}
	stack := &Stack{}
	queue := &Queue{}
	hashtable := &HashTable{capacity: 512, data: make([]*NodeHT, 512)}
	filename := "sqsh.txt"
	for scanner.Scan() {
		command := scanner.Text()
		Mute.Lock()
		response := processCommand(command, set, stack, queue, hashtable, filename)
		Mute.Unlock()
		_, err := conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Println("Ошибка при отправке ответа клиенту:", err)
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Ошибка при чтении команд от клиента:", scanner.Err())
	}
}
