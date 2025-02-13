package main

import (
	"SQLRecorder/mysql"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

func display_banner() {
	fmt.Println("   _____  ____  _      _____                        _           ")
	fmt.Println("  / ____|/ __ \\| |    |  __ \\     𝚋𝚢：𝙻𝚝𝚖𝚃𝚑𝚒𝚗𝚔     | |          ")
	fmt.Println(" | (___ | |  | | |    | |__) |___  ___ ___  _ __ __| | ___ _ __ ")
	fmt.Println("  \\___ \\| |  | | |    |  _  // _ \\/ __/ _ \\| '__/ _` |/ _ \\ '__|")
	fmt.Println("  ____) | |__| | |____| | \\ \\  __/ (_| (_) | | | (_| |  __/ |   ")
	fmt.Println(" |_____/ \\___\\_\\______|_|  \\_\\___|\\___\\___/|_|  \\__,_|\\___|_|   ")
	fmt.Println("                                                                ")
	fmt.Println("                                                                ")
}
func main() {
	display_banner()
	// 创建通道监听系统信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	var app = cli.App{
		Name:      "SQLRecorder",
		Usage:     "Create a proxy to record all passing SQL statements.",
		UsageText: "sqlrecorder command -s 127.0.0.1:3306 -p 127.0.0.1:43306",
		Commands: []*cli.Command{
			{
				Name:      "command",
				Aliases:   []string{"c"},
				Usage:     "Create a proxy and all SQL will be displayed in the command line window",
				UsageText: "sqlrecorder command -s 127.0.0.1:3306 -p 127.0.0.1:43306",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "server",
						Aliases:  []string{"s"},
						Usage:    "The address of the SQL server.",
						Required: false,
						Value:    "127.0.0.1:3306",
					},
					&cli.StringFlag{
						Name:     "proxy",
						Aliases:  []string{"p"},
						Usage:    "The address where the SQLRecorder agent wants to run the listening.",
						Required: false,
						Value:    "127.0.0.1:43306",
					},
				},
				Action: func(context *cli.Context) error {
					var server = context.String("server")
					var proxy = context.String("proxy")
					if server == "" || proxy == "" {
						return errors.New("Please enter the correct parameters")
					}
					sqlName := "mysql"
					switch sqlName {
					case "mysql":
						err := mysql.Recorder(server, proxy)
						return err
					default:
						return errors.New("Please enter the correct parameters")
					}
				},
			},
		},
	}
	// 隐藏光标
	fmt.Print("\033[?25l")
	go func() {
		<-sigs
		fmt.Println("\033[32m\nbye👋\033[0m")
		// 显示光标
		fmt.Print("\033[?25h")
		os.Exit(0)
	}()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m[Error] %v\n\033[0m", err.Error())
		fmt.Print("\033[?25h")
	}
}
