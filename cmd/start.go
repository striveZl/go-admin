package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"go-admin/internal/bootstrap"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func StartCmd() *cli.Command {

	return &cli.Command{
		Name:  "start",
		Usage: "start the server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "workdir",
				Aliases:     []string{"d"},
				Usage:       "Working directory",
				DefaultText: "configs",
				Value:       "configs",
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Runtime configuration files or directory (relative to workdir, multiple separated by commas)",
				DefaultText: "dev",
				Value:       "dev",
			},
			&cli.StringFlag{
				Name:    "static",
				Aliases: []string{"s"},
				Usage:   "Static files directory",
			},
			&cli.BoolFlag{
				Name:  "daemon",
				Usage: "Run as a daemon",
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			workdir := command.String("workdir")
			staticDir := command.String("static")
			configs := command.String("config")

			if command.Bool("daemon") {
				bin, err := filepath.Abs(os.Args[0]) //获取当前正在运行的可执行文件的「绝对路径」

				if err != nil {
					fmt.Printf("获取绝对路径失败: %s \n", err.Error())
					return err
				}

				args := []string{"start"}
				args = append(args, "-d", workdir)
				args = append(args, "-c", configs)
				args = append(args, "-s", staticDir)
				fmt.Printf("执行命令: %s %s \n", bin, strings.Join(args, " "))
				commands := exec.Command(bin, args...) //创建一个“即将启动的外部进程对象”

				//将标准输出和标准错误重定向到日志文件
				stdLogFile := fmt.Sprintf("%s.log", command.Name)

				file, err := os.OpenFile(stdLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					fmt.Printf("failed to open log file: %s \n", err.Error())
					return err
				}
				defer file.Close()

				commands.Stdout = file
				commands.Stderr = file

				err = commands.Start()
				if err != nil {
					fmt.Printf("failed to start daemon thread: %s \n", err.Error())
					return err
				}
				pid := 0
				if commands.Process != nil {
					pid = commands.Process.Pid
				}
				if pid <= 0 {
					return fmt.Errorf("failed to get daemon process pid")
				}

				serviceName := filepath.Base(bin)
				fmt.Printf("Service %s daemon thread started successfully\n", serviceName)
				if err := os.WriteFile(fmt.Sprintf("%s.pid", command.Name), []byte(fmt.Sprintf("%d", pid)), 0666); err != nil {
					return fmt.Errorf("failed to write pid file: %w", err)
				}
				fmt.Printf("service %s daemon thread started with pid %d \n", serviceName, pid)
				os.Exit(0)
			}

			//应用启动
			if err := bootstrap.Run(context.Background(), bootstrap.RunConfig{
				WorkDir:   workdir,
				Configs:   configs,
				StaticDir: staticDir,
			}); err != nil {
				return err
			}

			return nil
		},
	}
}
