package mysql

// 连接时传递的客户端配置
type clientConfig struct {
	isParse bool
	// 是否启用 Deprecate EOF
	deprecateEOF bool
}

// 暂时只解析 Deprecate EOF:
func (cf *clientConfig) parseClientConfig(data []byte) {
	// Deprecate EOF
	if data[7]&0x01 == 0 {
		cf.deprecateEOF = false
	} else {
		cf.deprecateEOF = true
	}
	cf.isParse = true
}

// 连接时传递的服务端配置
type serverConfig struct {
}

// 暂时掠过
func (cf *clientConfig) pareServerConfig(data []byte) {
	// pass
}
