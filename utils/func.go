package utils

import "net"

func SqlConnectTest(address string) error {
	serverConn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer serverConn.Close()
	return nil
}
