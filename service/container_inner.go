package service

import (
	"engine/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 同步类型的请求，函数上报结果
func (engine *Engine) FunctionEnd(ctx *gin.Context) {
	requestBody := model.FunctionEndRequest{}
	err := ctx.ShouldBindJSON(&requestBody)

	ctx.JSON(http.StatusOK, nil) // 直接返回200

	if err != nil {
		log.Error("FunctionEnd: ", err)
	} else {
		engine.containerService.FunctionEndHandler(requestBody)
	}
}

// 容器进程启动
func (engine *Engine) ContainerProcessStart(ctx *gin.Context) {
	id := ctx.Param("id")
	pidStr := ctx.Param("pid")
	timestampStr := ctx.Param("timestamp")

	ctx.JSON(http.StatusOK, nil) // 直接返回200

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Errorf("ContainerProcessEnd: id=%s, pid=%s, timestamp=%s, error=%v", id, pidStr, timestampStr, err)
		return
	}

	timestamp, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		log.Errorf("ContainerProcessEnd: id=%s, pid=%s, timestamp=%s, error=%v", id, pidStr, timestampStr, err)
		return
	}

	engine.containerService.ContainerProcessStartHandler(id, pid, timestamp)
}

// 容器进程结束，记录一些基本信息
func (engine *Engine) ContainerProcessEnd(ctx *gin.Context) {
	id := ctx.Param("id")
	timestamp := ctx.Param("timestamp")

	ctx.JSON(http.StatusOK, nil) // 直接返回200

	if t, err := strconv.ParseInt(timestamp, 10, 64); err != nil {
		log.Errorf("ContainerProcessEnd: id=%s, timestamp=%s, error=%v", id, timestamp, err)
	} else {
		go engine.containerService.ContainerProcessEndHandler(id, t)
	}
}
