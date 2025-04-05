package postgresql

const (
	maxPacketSize = 1<<24 - 1
)

//协议标识符见https://www.postgresql.org/docs/13/protocol-message-formats.html

const (
	iReady         byte = 0x5a
	iERR           byte = 0x45
	iComplete      byte = 0x43
	iDataRow       byte = 0x44
	iRowDesc       byte = 0x54
	iAuth          byte = 0x52
	iParseComplete byte = 0x31
)

const (
	comParse    byte = 0x50
	comQuery    byte = 0x51
	comBind     byte = 0x42
	comDescribe byte = 0x44
	comExecute  byte = 0x45
	comSync     byte = 0x53
	comPassword byte = 0x70
)
