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
var zygote bool
var parallel bool
var runtimeName string
var templateName string

func init() {
	flag.StringVar(&configPath, "c", "", "")
	flag.IntVar(&count, "n", 1, "")
	flag.BoolVar(&zygote, "zygote", false, "")
	flag.BoolVar(&parallel, "p", false, "")
	flag.StringVar(&runtimeName, "runtime", "python3.7", "")
	flag.StringVar(&templateName, "template", "echo", "")
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

	time.Sleep(3 * time.Second)
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
