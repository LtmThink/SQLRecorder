package postgresql

import (
	. "SQLRecorder/utils"
	"bytes"
	"fmt"
	"strings"
	"time"
)

type messages struct {
	messageList []message
	num         int
	isReady     bool
}

func (ms *messages) addMessageList(data []byte) {
	tmpPacket := message{
		packets: [][]byte{data},
		num:     0,
	}
	ms.messageList = append(ms.messageList, tmpPacket)
	ms.isReady = false
	ms.num++
}

func (ms *messages) addMessage(data []byte, num int) {
	packetNum := &ms.messageList[num].num
	ms.messageList[num].packets = append(ms.messageList[num].packets, data)
	*packetNum++
}

func (ms *messages) writeToTerminal() {
	comType := ms.messageList[ms.num].comType
	resType := ms.messageList[ms.num].resType
	if ms.isReady == false && comType != 0 && resType != 0 {
		ms.messageList[ms.num].print()
	}
	ms.isReady = true
}

// 请求处理
func (ms *messages) queryRecord(data []byte) {
	ms.messageList[ms.num].comType = comQuery
	query := data[5:]
	query = query[:len(query)-1]
	ms.messageList[ms.num].command = string(query)
}

func (ms *messages) parseRecord(data []byte) {
	ms.messageList[ms.num].comType = comParse
	msg := data[5:]
	parts := bytes.Split(msg, []byte{0x00})
	ms.messageList[ms.num].command = fmt.Sprintf("Parse '%s' as '%s'", string(parts[1]), string(parts[0]))
}

func (ms *messages) bindRecord(data []byte) {
	ms.messageList[ms.num].comType = comBind
	// 避免portal为空的情况我们自己给他加一个字节的值
	msg := []byte{0x65}
	msg = append(msg, data[5:]...)
	parts := bytes.Split(msg, []byte{0x00})
	ms.messageList[ms.num].command = fmt.Sprintf("Bind '%s' to Execute", string(parts[1]))
}

// 消息处理
func (ms *messages) errRecord(data []byte) {
	ms.messageList[ms.num].resType = iERR
	errorMessage := data[27:]
	// 以0X00的位置为间隔,前面部分为error信息
	errorMessage = errorMessage[:bytes.Index(errorMessage, []byte{0x00})]

	ms.messageList[ms.num].result += LightRed(fmt.Sprintf("ERROR : %s ;", errorMessage))

}
func (ms *messages) completeRecord(data []byte) {
	ms.messageList[ms.num].resType = iComplete
	completeMessage := string(data[5:])
	completeMessage = completeMessage[:len(completeMessage)-1]
	if strings.Contains(completeMessage, "SELECT") {
		ms.messageList[ms.num].result += LightYellow(fmt.Sprintf("Query OK , %d columns, %d rows in set ;", ms.messageList[ms.num].columnNum, ms.messageList[ms.num].rowNum))
	} else {
		ms.messageList[ms.num].result += LightYellow(fmt.Sprintf("OK , %s ;", completeMessage))
	}
}

type message struct {
	packets   [][]byte
	num       int
	comType   byte
	command   string
	resType   byte
	result    string
	columnNum int
	rowNum    int
}

func (m *message) print() {
	timestamp := time.Now().Format("15:04:05")
	comTypeString := ""
	switch m.comType {
	case comQuery:
		comTypeString = "Query"
	case comParse:
		comTypeString = "Parse"
	case comBind:
		comTypeString = "Bind"
	default:
		//未定义comType无法记录
		return
	}
	fmt.Printf("\n%s -> [%s] `%s` -> `%s`", Green(comTypeString), Yellow(timestamp), LightGreen(m.command), m.result)
}

func (m *message) processTabular() {
}
