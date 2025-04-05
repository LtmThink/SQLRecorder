package postgresql

// 连接时传递的客户端配置
type clientConfig struct {
	isParse bool
}

// 暂时掠过
func (cf *clientConfig) parseClientConfig(data []byte) {
	// 解析客户端配置
	// pass
}

// 连接时传递的服务端配置
type serverConfig struct {
}

// 暂时掠过
func (cf *clientConfig) pareServerConfig(data []byte) {
	// pass
}
