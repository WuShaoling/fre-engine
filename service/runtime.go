package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (engine *Engine) ListRuntime(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, engine.templateService.List())
}

func (engine *Engine) DumpRuntime(ctx *gin.Context) {
	if err := engine.templateService.Dump(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": "ok"})
	}
}
