package mysql

import (
	"fmt"
	"log"
	"net"
)

func Recorder(server string, proxy string) (err error) {
	err = sqlConnectTest(server)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", proxy)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println(Yellow("Please make sure the mysql server is running on " + server + " ..."))
	fmt.Println(Green("MySQL Proxy is listening on " + proxy + " ..."))

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

func sqlConnectTest(address string) error {
	mysqlConn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer mysqlConn.Close()
	return nil
}

func handleClientConnection(clientConn net.Conn, server string) {
	// 连接到实际的 MySQL 服务器
	mysqlConn, err := net.Dial("tcp", server)
	if err != nil {
		log.Println("[Proxy] Error: %v", err)
		return
	}
	// 启动两个 goroutine，一个处理客户端请求，另一个转发 MySQL 响应
	go forwardToMySQL(clientConn, mysqlConn)
	go forwardToClient(mysqlConn, clientConn)
}

func forwardToMySQL(clientConn, mysqlConn net.Conn) {
	defer mysqlConn.Close() // 确保在所有数据转发完成后再关闭 mysqlConn
	clientRecorder := newConn(clientConn)
	for {
		n, err := clientRecorder.recordClientPacket()
		if err != nil {
			return
		}
		// 转发客户端的数据到 MySQL 服务器
		_, err = mysqlConn.Write(n)
		if err != nil {
			return
		}
	}
}

func forwardToClient(mysqlConn, clientConn net.Conn) {
	defer clientConn.Close() // 确保在所有数据转发完成后再关闭 clientConn
	serverRecorder := newConn(mysqlConn)
	for {
		n, err := serverRecorder.recordServerPacket()
		if err != nil {
			return
		}
		// 转发 MySQL 响应到客户端
		_, err = clientConn.Write(n)
		if err != nil {
			return
		}
	}
}
