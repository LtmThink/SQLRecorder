package mysql

import (
	"fmt"
	"net"
	"time"
)

type conn struct {
	buf      buffer
	sequence uint8
}

func newConn(nc net.Conn) conn {
	fg := newBuffer(nc)
	return conn{fg, 0}
}

func (c *conn) readPacket() ([]byte, []byte, error) {
	var compData []byte
	var prevData []byte
	for {
		// read packet header
		data, err := c.buf.readNext(4)
		compData = append(compData, data...)
		if err != nil {
			return nil, nil, ErrInvalidConn
		}

		// packet length [24 bit]
		pktLen := int(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16)

		//// check packet sync [8 bit]
		//fmt.Println(data)
		//if data[3] != c.sequence {
		//	if data[3] > c.sequence {
		//		return nil, nil, ErrPktSyncMul
		//	}
		//	return nil, nil, ErrPktSync
		//}
		//c.sequence++

		// packets with length 0 terminate a previous packet which is a
		// multiple of (2^24)-1 bytes long
		if pktLen == 0 {
			// there was no previous packet
			if prevData == nil {
				return nil, nil, ErrInvalidConn
			}

			return prevData, compData, nil
		}

		// read packet body [pktLen bytes]
		data, err = c.buf.readNext(pktLen)
		compData = append(compData, data...)
		if err != nil {
			return nil, compData, ErrInvalidConn
		}

		// return data if this was the last packet
		if pktLen < maxPacketSize {
			// zero allocations for non-split packets
			if prevData == nil {
				return data, compData, nil
			}

			return append(prevData, data...), compData, nil
		}

		prevData = append(prevData, data...)
	}
}
func (c *conn) recordClientPacket() ([]byte, error) {
	data, compData, err := c.readPacket()
	if err != nil {
		return nil, err
	}
	if data[0] == comQuery {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("\n%s -> [%s] `%s`", Green("Query"), Yellow(timestamp), LightGreen(string(data[1:])))
	}
	return compData, nil
}
func (c *conn) recordServerPacket() ([]byte, error) {
	data, compData, err := c.readPacket()
	if err != nil {
		return nil, err
	}
	if data[0] == iERR {
		error_code := int(data[2])<<8 | int(data[1])
		sqlstate := string(data[4:9])
		error_message := string(data[9:])
		fmt.Printf(" -> `%s`", LightRed("ERROR "+fmt.Sprintf("%d", error_code)+" ("+sqlstate+"): "+error_message))
	}
	return compData, nil
}
