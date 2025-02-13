package mysql

import (
	"net"
	"testing"
)

func TestRecorder(t *testing.T) {
	type args struct {
		server string
		proxy  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "right proxy",
			args:    args{"127.0.0.1:3306", "127.0.0.1:43036"},
			wantErr: true,
		},
		{
			name:    "error proxy",
			args:    args{"127.0.0.1:8588", "127.0.0.1:43036"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Recorder(tt.args.server, tt.args.proxy); (err != nil) != tt.wantErr {
				t.Errorf("Recorder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_forwardToClient(t *testing.T) {
	type args struct {
		mysqlConn  net.Conn
		clientConn net.Conn
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forwardToClient(tt.args.mysqlConn, tt.args.clientConn)
		})
	}
}

func Test_forwardToMySQL(t *testing.T) {
	type args struct {
		clientConn net.Conn
		mysqlConn  net.Conn
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forwardToMySQL(tt.args.clientConn, tt.args.mysqlConn)
		})
	}
}

func Test_handleClientConnection(t *testing.T) {
	type args struct {
		clientConn net.Conn
		server     string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleClientConnection(tt.args.clientConn, tt.args.server)
		})
	}
}

func Test_sqlConnectTest(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sqlConnectTest(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("sqlConnectTest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
