package main

import (
	"context"
	"flag"
	"github.com/unchainese/unchain/global"
	"github.com/unchainese/unchain/server"
	"os"
	"os/signal"
	"time"
)

var configFilePath, installMode, action string

func main() {
	flag.StringVar(&action, "action", "run", "动作参数,可选值: run, install,uninstall,info,run")
	flag.StringVar(&configFilePath, "config", "config.toml", "配置文件路径")
	flag.StringVar(&installMode, "mode", "single", "安装命令的模式参数")
	flag.Parse()

	c := global.Cfg(configFilePath) //using default config.toml file
	fd := global.SetupLogger(c)
	defer fd.Close()

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	app := node.NewApp(c, stop)
	app.PushNode()                 //register node info to the manager server
	app.PrintVLESSConnectionURLS() //for standalone node
	go app.Run()
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Shutdown(ctx)

}
