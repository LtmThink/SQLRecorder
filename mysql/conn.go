package mysql

import (
	"fmt"
	"net"
)

type conn struct {
	buf          buffer
	messages     *messages
	clientConfig *clientConfig
}

func newConn(nc net.Conn, p *messages, cf *clientConfig) conn {
	fg := newBuffer(nc)
	return conn{fg, p, cf}
}

func (c *conn) readPacket() ([]byte, []byte, error) {
	var compData []byte
	var prevData []byte
	for {
		// read message header
		data, err := c.buf.readNext(4)
		compData = append(compData, data...)
		if err != nil {
			return nil, nil, ErrInvalidConn
		}

		// message length [24 bit]
		pktLen := int(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16)

		//// check message sync [8 bit]
		//fmt.Println(packets)
		//if packets[3] != c.sequence {
		//	if packets[3] > c.sequence {
		//		return nil, nil, ErrPktSyncMul
		//	}
		//	return nil, nil, ErrPktSync
		//}
		//c.sequence++

		// messages with length 0 terminate a previous message which is a
		// multiple of (2^24)-1 bytes long
		if pktLen == 0 {
			// there was no previous message
			if prevData == nil {
				return nil, nil, ErrInvalidConn
			}

			return prevData, compData, nil
		}

		// read message body [pktLen bytes]
		data, err = c.buf.readNext(pktLen)
		compData = append(compData, data...)
		if err != nil {
			return nil, compData, ErrInvalidConn
		}

		// return packets if this was the last message
		if pktLen < maxPacketSize {
			// zero allocations for non-split messages
			if prevData == nil {
				return data, compData, nil
			}

			return append(prevData, data...), compData, nil
		}

		prevData = append(prevData, data...)
	}
}

func (c *conn) recordPacket() ([]byte, error) {
	ms := c.messages

	_, data, err := c.readPacket()
	// 除了读取数据失败其他情况都不返回error
	if err != nil {
		return nil, err
	}
	// 解析配置
	if !c.clientConfig.isParse {
		if data[3] == 0x01 {
			c.clientConfig.parseClientConfig(data)
		}
	}

	// 客户端查询解析
	if data[3] == 0x00 {
		// 如果是第一个message
		if ms.num == -1 {
			ms.addMessageList(data)
			return data, nil
		}
		// 如果不超过256个packet或者刚好第256个packet是EOF
		if ms.messageList[ms.num].num < 255 || (ms.messageList[ms.num].num >= 255 && ms.messageList[ms.num].packets[ms.messageList[ms.num].num][4] == iEOF) {
			ms.addMessageList(data)
			switch data[4] {
			case comQuery:
				ms.messageList[ms.num].command = string(data[5:])
				ms.messageList[ms.num].comType = comQuery
			case comInitDB:
				ms.messageList[ms.num].command = "use " + string(data[5:])
				ms.messageList[ms.num].comType = comInitDB
			default:
				//pass
			}
			return data, nil
		}
		return data, nil
	}

	// 服务端响应解析
	// 排除登录阶段（第一个message）
	num := ms.num
	if num > 0 {
		switch data[4] {
		case iERR:
			// 记录到packetList
			errorCode := int(data[6])<<8 | int(data[5])
			sqlState := string(data[8:13])
			errorMessage := string(data[13:])
			ms.messageList[num].result = fmt.Sprintf("ERROR %d (%s): %s", errorCode, sqlState, errorMessage)
			ms.messageList[num].resType = iERR
			// 记录信息到终端
			ms.messageList[num].writeToTerminal()
		case iOK:
			// 记录到packetList
			affectedRows := int(data[5])
			ms.messageList[num].result = fmt.Sprintf("OK, %d rows affected", affectedRows)
			ms.messageList[num].resType = iOK
			// 记录信息到终端
			ms.messageList[num].writeToTerminal()
		case iEOF:
			//两种情况处理
			//field中间没有eof处理或中间有eof但是是第二个eof了
			if c.clientConfig.deprecateEOF || ms.messageList[num].eofNum == 1 {
				ms.messageList[num].processTabular(c.clientConfig.deprecateEOF)
				ms.messageList[num].writeToTerminal()
			}
			ms.messageList[num].eofNum++

		default:
			//pass
		}
	}
	ms.addMessage(data, ms.num)
	return data, nil
}
