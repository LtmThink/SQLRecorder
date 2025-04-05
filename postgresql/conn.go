package postgresql

import (
	. "SQLRecorder/buffer"
	. "SQLRecorder/utils"
	"net"
)

type conn struct {
	buf          Buffer
	messages     *messages
	clientConfig *clientConfig
	isClient     bool
}

func newConn(nc net.Conn, p *messages, cf *clientConfig, isClient bool) conn {
	fg := NewBuffer(nc)
	return conn{fg, p, cf, isClient}
}

func (c *conn) readPacket() ([]byte, []byte, error) {
	var compData []byte
	var prevData []byte
	for {
		// read message header
		data, err := c.buf.ReadNext(5)
		compData = append(compData, data...)
		if err != nil {
			return nil, nil, ErrInvalidConn
		}

		// message length [32 bit]
		pktLen := int(uint32(data[4]) | uint32(data[3])<<8 | uint32(data[2])<<16 | uint32(data[1])<<24)

		if pktLen == 0 {
			// there was no previous message
			if prevData == nil {
				return nil, nil, ErrInvalidConn
			}

			return prevData, compData, nil
		}

		// read message body [pktLen bytes]
		data, err = c.buf.ReadNext(pktLen - 4)
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
func (c *conn) readStartupPacket() ([]byte, error) {
	for {
		// read message header
		data, err := c.buf.ReadAll()
		if err != nil {
			return nil, ErrInvalidConn
		}
		return data, nil
	}
}
func (c *conn) recordPacket() ([]byte, error) {
	ms := c.messages
	var data []byte
	var err error
	// 每个连接的前几个包是启动包，比较特殊跳过处理
	if c.clientConfig.isParse {
		_, data, err = c.readPacket()
		if err != nil {
			return nil, err
		}
		if ms.isReady {
			ms.addMessageList(data)
		} else {
			ms.addMessage(data, ms.num)
		}
		if c.isClient {
			// 处理请求
			switch data[0] {
			case comQuery:
				ms.queryRecord(data)
			case comParse:
				ms.parseRecord(data)
			case comBind:
				ms.bindRecord(data)
			}
		} else {
			// 处理响应
			switch data[0] {
			case iComplete:
				ms.completeRecord(data)
			case iDataRow:
				ms.messageList[ms.num].rowNum++
			case iRowDesc:
				ms.messageList[ms.num].columnNum = int(data[5])<<8 | int(data[6])
			case iReady:
				ms.writeToTerminal()
			case iERR:
				ms.errRecord(data)
			case iParseComplete:
				// 解析完成
				ms.messageList[ms.num].resType = iParseComplete
				ms.messageList[ms.num].result += LightYellow("Parse Complete ;")
			}
		}
	} else {
		data, err = c.readStartupPacket()
		if err != nil {
			return nil, err
		}
		if data[0] == comPassword {
			// 跳过前面的几个包，通过也就是完成了初始化解析
			c.clientConfig.isParse = true
		}
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}
