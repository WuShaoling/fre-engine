package service

import (
	"engine/template"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (engine *Engine) CreateTemplate(ctx *gin.Context) {
	requestBody := &template.Template{}
	if err := ctx.ShouldBindJSON(requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数不完整"})
		return
	}

	if err := engine.templateService.Create(requestBody); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"data": "ok"})
	}
}

func (engine *Engine) ListTemplate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, engine.templateService.List())
}

func (engine *Engine) DeleteTemplate(ctx *gin.Context) {
	name := ctx.Param("name")
	if err := engine.templateService.Delete(name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": "ok"})
	}
}

func (engine *Engine) DumpTemplate(ctx *gin.Context) {
	if err := engine.templateService.Dump(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": "ok"})
	}
}
