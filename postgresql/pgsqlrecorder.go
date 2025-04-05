package postgresql

import (
	. "SQLRecorder/utils"
	"fmt"
	"log"
	"net"
)

func Recorder(server string, proxy string) (err error) {
	err = SqlConnectTest(server)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", proxy)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println(Yellow("Please make sure the PostgreSQL server is running on " + server + " ..."))
	fmt.Println(Green("PostgreSQL Proxy is listening on " + proxy + " ..."))
	fmt.Println(Red("Press Ctrl-C to quit ..."))

	for {
		// 等待客户端连接
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting client connection:", err)
			continue
		}
		// 启动一个新的协程来处理该客户端连接
		go handleClientConnection(clientConn, server)
	}
}

func handleClientConnection(clientConn net.Conn, server string) {
	// 连接到实际的 MySQL 服务器
	serverConn, err := net.Dial("tcp", server)
	if err != nil {
		log.Println("[Proxy] Error: %v", err)
		return
	}
	// 初始化packet记录器
	p := messages{num: -1, isReady: true}
	// 初始化客户端配置
	cf := clientConfig{isParse: false}
	// 启动两个 goroutine，一个处理客户端请求，另一个转发 MySQL 响应
	clientRecorder := newConn(clientConn, &p, &cf, true)
	serverRecorder := newConn(serverConn, &p, &cf, false)
	go sendToReceive(serverConn, clientRecorder)
	go sendToReceive(clientConn, serverRecorder)
}

func sendToReceive(receiveConn net.Conn, c conn) {
	defer receiveConn.Close() // 确保在所有数据转发完成后再关闭
	for {
		n, err := c.recordPacket()
		if err != nil {
			return
		}
		_, err = receiveConn.Write(n)
		if err != nil {
			return
		}
	}
}
