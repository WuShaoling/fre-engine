package container

import (
	"engine/runtime"
	"engine/template"
	"errors"
)

// 通过 zygote 启动容器进程

func (service *Service) newContainerProcessByZygote(r *runtime.Runtime, t *template.Template, container *Container) error {
	// 构建参数
	functionExecContext := service.buildFunctionExecContext(t, container)
	if functionExecContext == "" {
		return errors.New("BuildFunctionExecContextError")
	}

	return service.zygoteService.NewContainerByZygoteProcess(r.Name, t.Name, functionExecContext)
}
