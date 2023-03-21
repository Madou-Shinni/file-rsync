package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

//go:generate go run main.go --ip 192.168.110.94 --src /cygdrive/D/go-project/frisbee-officer-backend-GVA/server/uploads/file/ --dist /root/rsfile

func main() {
	// 本地上传目录和远程同步目录
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "ip",
				Aliases:  []string{"I"},
				Usage:    "远程主机的ip",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "src",
				Aliases:  []string{"S"},
				Usage:    "需要同步文件的目录",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "dist",
				Aliases: []string{"D"},
				Usage:   "需要同步文件的目标目录，default:src目录",
			},
			&cli.StringFlag{
				Name:        "module",
				Aliases:     []string{"M"},
				DefaultText: "test",
				Usage:       "rsync的模组(需要提前在配置文件中创建)，default:test",
			},
		},
		Action: func(c *cli.Context) error {
			ip := c.String("ip")
			path := c.String("src")
			dist := c.String("dist")
			module := c.String("module")
			if dist == "" {
				dist = path
			}
			if module == "" {
				module = "test"
			}
			// 定时任务
			AsyncScheduler(ip, path, dist, module)

			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			select {
			case <-signals:
				// 释放资源
				fmt.Println("[file-rsync] 程序关闭，释放资源")
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Synchronization 同步文件
func Synchronization(ip, src, dist, module string) {
	// rsync 是一个常用的 Linux 应用程序，用于文件同步
	// rsync -avz --progress /cygdrive/D/go-project/frisbee-officer-backend-GVA/server/uploads/file/ 192.168.110.94::test/rsfile
	//cmd := exec.Command("rsync", "-avz", "--progress", src, ip+"::test"+dist)
	cmd := exec.Command("rsync", "-avz", "--progress", src, ip+"::"+module+dist)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("rsync error: %v\n%s", err, output)
		// 退出程序
		os.Exit(os.Getgid())
	}
}

// pathNotExistCreate 目录不存在就创建目录
func pathNotExistCreate(folderPath string) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 文件夹不存在，创建该文件夹
		err = os.MkdirAll(folderPath, 0755)
		if err != nil {
			fmt.Println("创建文件夹失败：", err)
		} else {
			fmt.Println("文件夹创建成功！")
		}
	} else {
		// 文件夹已经存在，不需要进行任何操作
		fmt.Println("文件夹已经存在，无需创建。")
	}
}

// AsyncScheduler 定时任务
// params: ip 远程主机的ip地址 src 源目录 dist 目标目录 module rsync的模组
func AsyncScheduler(ip string, src string, dist, module string) {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	s := gocron.NewScheduler(timezone)
	// 每隔一分钟执行一次
	s.Every(1).Minute().Do(func() {
		Synchronization(ip, src, dist, module)
	})
	s.StartAsync()
}
