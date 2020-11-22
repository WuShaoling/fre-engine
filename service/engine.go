package service

import (
	"engine/container"
	"engine/runtime"
	"engine/template"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	functionResultWaitChanMap  map[string]chan gin.H
	functionResultCallbackChan chan gin.H
	containerExitCallbackChan  chan gin.H

	runtimeService   *runtime.Service
	templateService  *template.Service
	containerService *container.Service
}

func NewEngine() *Engine {
	functionResultWaitChanMap := make(map[string]chan gin.H)
	functionResultCallbackChan := make(chan gin.H, 16)
	containerExitCallbackChan := make(chan gin.H, 16)

	runtimeService := runtime.NewRuntimeService()
	templateService := template.NewTemplateService()
	containerService := container.NewContainerService(containerExitCallbackChan, functionResultCallbackChan,
		runtimeService.List(), templateService.List())

	engine := &Engine{
		functionResultWaitChanMap:  functionResultWaitChanMap,
		containerExitCallbackChan:  containerExitCallbackChan,
		functionResultCallbackChan: functionResultCallbackChan,

		runtimeService:   runtimeService,
		templateService:  templateService,
		containerService: containerService,
	}
	engine.startContainerExitCallbackHandler()
	return engine
}

func (engine *Engine) startContainerExitCallbackHandler() {
	go func() {
		var response gin.H
		for {
			select {
			case response = <-engine.containerExitCallbackChan:
				if _, ok := engine.functionResultWaitChanMap[response["requestId"].(string)]; ok {
					// 容器已经退出，chan 还在，说明函数执行结果未上报
					response["error"] = "FunctionExitButNotReportResult"
				}
			case response = <-engine.functionResultCallbackChan:
			}
			requestId := response["requestId"].(string)
			if c, ok := engine.functionResultWaitChanMap[requestId]; ok {
				c <- response
				close(c)
				delete(engine.functionResultWaitChanMap, requestId)
			}
		}
	}()
}
