package mysql

import "github.com/gookit/color"

const (
	maxPacketSize = 1<<24 - 1
)

const (
	iOK           byte = 0x00
	iAuthMoreData byte = 0x01
	iLocalInFile  byte = 0xfb
	iEOF          byte = 0xfe
	iERR          byte = 0xff
)

const (
	comQuit byte = iota + 1
	comInitDB
	comQuery
	comFieldList

	comCreateDB
	comDropDB
	comRefresh
	comShutdown
	comStatistics
	comProcessInfo
	comConnect
	comProcessKill
	comDebug
	comPing
	comTime
	comDelayedInsert
	comChangeUser
	comBinlogDump
	comTableDump
	comConnectOut
	comRegisterSlave
	comStmtPrepare
	comStmtExecute
	comStmtSendLongData
	comStmtClose
	comStmtReset
	comSetOption
	comStmtFetch
)

var (
	// color
	Green       = color.Green.Render
	LightGreen  = color.LightGreen.Render
	Yellow      = color.Yellow.Render
	LightYellow = color.LightYellow.Render
	Red         = color.Red.Render
	LightRed    = color.LightRed.Render
)
