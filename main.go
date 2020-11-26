package main

import (
	"engine/config"
	"engine/container"
	"engine/runtime"
	"engine/template"
	"flag"
	"os"
	"time"
)

var configPath string
var count int
var parallel bool
var zygote bool
var runtimeName string
var templateName string

func init() {
	flag.IntVar(&count, "n", 1, "创建的数量")
	flag.BoolVar(&parallel, "p", false, "并发启动")
	flag.BoolVar(&zygote, "zygote", false, "并发启动")
	flag.StringVar(&runtimeName, "runtime", "python3.7", "并发启动")
	flag.StringVar(&templateName, "template", "normal", "并发启动")
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "exec" { // 容器进程
		container.Exec()
		return
	}

	flag.Parse()
	config.InitSysConfig(configPath)

	runtimeService := runtime.NewRuntimeService()
	templateService := template.NewTemplateService()
	containerService := container.NewContainerService(runtimeService.List(), templateService.List())
	r := runtimeService.Get(runtimeName)
	t := templateService.Get(templateName)

	if parallel {
		for i := 0; i < count; i++ {
			go func(id int) {
				_, _ = containerService.Create(id, r, t, zygote, map[string]interface{}{})
			}(i)
		}
	} else {
		for i := 0; i < count; i++ {
			_, _ = containerService.Create(i, r, t, zygote, map[string]interface{}{})
		}
	}
	time.Sleep(time.Hour)
}
