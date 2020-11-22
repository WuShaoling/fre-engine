package api

import (
	"engine/service"
	"github.com/gin-gonic/gin"
)

func SetContainerRouter(engine *service.Engine, r *gin.Engine) {
	// 对外暴露的
	runtimeGroup := r.Group("/api/runtime")
	runtimeGroup.GET("/", engine.ListRuntime)
	runtimeGroup.PUT("/dump", engine.DumpRuntime)

	templateGroup := r.Group("/api/template")
	templateGroup.GET("/", engine.ListTemplate)
	templateGroup.POST("/", engine.CreateTemplate)
	templateGroup.DELETE("/:name", engine.DeleteTemplate)
	templateGroup.PUT("/dump", engine.DumpTemplate)

	containerGroup := r.Group("/api/container")
	containerGroup.GET("/", engine.ListContainer)
	containerGroup.POST("/:template", engine.CreateContainer)
	containerGroup.PUT("/dump", engine.DumpContainer)

	// 内部使用的
	innerGroup := r.Group("/inner")
	innerGroup.PUT("/function/end", engine.FunctionEnd)
	innerGroup.PUT("/process/run/:id/:timestamp/:pid", engine.ContainerProcessStart) // zygote使用
	innerGroup.PUT("/process/end/:id/:timestamp", engine.ContainerProcessEnd)        // zygote使用
}
