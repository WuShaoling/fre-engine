package service

import (
	"engine/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (engine *Engine) CreateContainer(ctx *gin.Context) {

	templateName := ctx.Param("template")
	functionParam := make(map[string]interface{})
	sync := ctx.DefaultQuery("sync", "false")

	if err := ctx.ShouldBindJSON(&functionParam); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "BadParameter"})
		return
	}

	template, ok := engine.templateService.Get(templateName)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TemplateNotExit"})
		return
	}

	runtime, ok := engine.runtimeService.Get(template.Runtime)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "RuntimeNotExit"})
		return
	}

	requestId := util.UniqueId()
	id, err := engine.containerService.Create(requestId, runtime, template, functionParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sync != "true" {
		c := make(chan gin.H, 1)
		engine.functionResultWaitChanMap[requestId] = c
		response := <-c
		fmt.Println(response)
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"id": id})
	}
}

func (engine *Engine) ListContainer(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, engine.containerService.List())
}

func (engine *Engine) StopContainer(ctx *gin.Context) {
	// TODO 把容器进程直接杀掉
	ctx.JSON(http.StatusOK, "ok")
}

func (engine *Engine) DumpContainer(ctx *gin.Context) {
	if err := engine.containerService.Dump(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": "ok"})
	}
}
