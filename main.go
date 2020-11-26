package main

import (
	"engine/api"
	"engine/config"
	"engine/container"
	"engine/service"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "config path")
}

func initLog() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true                    // 显示完整时间
	customFormatter.TimestampFormat = "2006-01-02 15:04:05" // 时间格式
	customFormatter.DisableTimestamp = false                // 禁止显示时间
	customFormatter.DisableColors = false                   // 禁止颜色显示
	logrus.SetFormatter(customFormatter)
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "exec" { // 容器进程
		container.Exec()
		return
	}

	initLog()
	flag.Parse()

	config.InitSysConfig(configPath)
	gin.SetMode(gin.ReleaseMode)

	freEngine := service.NewEngine()
	r := gin.Default()
	api.SetContainerRouter(freEngine, r)

	log.Println("server listen on :" + config.SysConfigInstance.ServePort)
	_ = r.Run(":" + config.SysConfigInstance.ServePort)
}
