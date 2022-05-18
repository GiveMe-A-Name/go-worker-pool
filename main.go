package main

import (
	"fmt"
	"net"
	"strings"
)

func worker(p chan int, results chan string) {
	for port := range p {
		address := fmt.Sprintf("localhost:%d", port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- fmt.Sprintf("1端口关闭%d\n", port)
			continue
		}
		conn.Close()
		results <- fmt.Sprintf("2端口打开了!!!%d\n", port)
	}
}

func main() {
	ports := make(chan int, 100)
	results := make(chan string)
	var opened = make([]string, 0, 500)
	var closed = make([]string, 0, 500)

	for i := 0; i < cap(ports); i++ {
		// 生成worker
		go worker(ports, results)
	}

	go func() {
		// 生成任务的goroutine
		for i := 1; i < 1024; i++ {
			ports <- i
		}
	}()
	// 主程序做为主goroutine来消费results
	for i := 1; i < 1024; i++ {
		result := <-results
		fmt.Println(result)
		if strings.Contains(result, "关闭") {
			closed = append(closed, result)
		} else if strings.Contains(result, "打开") {
			opened = append(opened, result)
		}
	}

	close(ports)
	close(results)
	fmt.Println(opened)
	fmt.Println(closed)
}
