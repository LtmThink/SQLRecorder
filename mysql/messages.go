package mysql

import (
	. "SQLRecorder/utils"
	"fmt"
	"time"
)

type messages struct {
	messageList []message
	num         int
}

func (ms *messages) addMessageList(data []byte) {
	tmpPacket := message{
		packets:  [][]byte{data},
		num:      0,
		eofNum:   0,
		printNum: 0,
	}
	ms.messageList = append(ms.messageList, tmpPacket)
	ms.num++
}
func (ms *messages) addMessage(data []byte, num int) {
	packetNum := &ms.messageList[num].num
	ms.messageList[num].packets = append(ms.messageList[num].packets, data)
	*packetNum++
}

type message struct {
	packets  [][]byte
	num      int
	eofNum   int
	comType  byte
	command  string
	resType  byte
	result   string
	printNum int //打印次数
}

func (m *message) writeToTerminal() {
	timestamp := time.Now().Format("15:04:05")
	comTypeString := ""
	switch m.comType {
	case comQuery:
		comTypeString = "Query"
	case comInitDB:
		comTypeString = "InitDB"
	default:
		//未定义comType无法记录
		return
	}
	if m.resType != iERR {
		if m.printNum == 0 {
			fmt.Printf("\n%s -> [%s] `%s` -> `%s`", Green(comTypeString), Yellow(timestamp), LightGreen(m.command), LightYellow(m.result))
		} else if m.printNum > 0 {
			fmt.Printf("\n-> `%s`", LightYellow(m.result))
		}
	} else {
		if m.printNum == 0 {
			fmt.Printf("\n%s -> [%s] `%s` -> `%s`", Green(comTypeString), Yellow(timestamp), LightGreen(m.command), LightRed(m.result))
		} else if m.printNum > 0 {
			fmt.Printf("\n-> `%s`", LightRed(m.result))
		}
	}
	m.printNum++

}

func (m *message) processTabular(deprecateEOF bool) {
	if m.num >= 255 {
		// packet包数超过255个，不处理
		m.result = fmt.Sprintf("OK, too many columns and rows")
		m.resType = iEOF
		return
	}
	fieldNumber := 1 + m.printNum
	fieldCount := int(m.packets[fieldNumber][4])
	rawNumber := fieldNumber + 1 + fieldCount
	if !deprecateEOF {
		rawNumber++
	}
	rawCount := m.num - rawNumber + 1
	m.result = fmt.Sprintf("OK, %d columns, %d rows in set", fieldCount, rawCount)
	m.resType = iEOF
	// 获取具体查询表格结果,留待以后再写吧
	// 获取字段

	// 获取行
}
