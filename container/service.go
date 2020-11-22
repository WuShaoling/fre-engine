package container

import (
	"engine/config"
	"engine/runtime"
	"engine/template"
	"engine/util"
	"errors"
	"github.com/gin-gonic/gin"
	"path"
	"time"
)

type Service struct {
	onContainerExitCallback  chan gin.H
	onFunctionResultCallback chan gin.H

	dataMap map[string]*Container

	fsService     *FsService
	cgroupService *CgroupService
	zygoteService *ZygoteService
}

func NewContainerService(containerExitCallback chan gin.H, functionResultCallback chan gin.H,
	runtimeSet map[string]*runtime.Runtime, templateSet map[string]*template.Template) *Service {
	service := &Service{
		onContainerExitCallback:  containerExitCallback,
		onFunctionResultCallback: functionResultCallback,

		dataMap:   make(map[string]*Container),
		fsService: NewFsService(),

		cgroupService: NewCgroupService(),
		zygoteService: NewZygoteService(runtimeSet, templateSet),
	}
	util.LoadJsonDataFromFile(service.getDataFilePath(), &service.dataMap)
	return service
}

func (service *Service) Create(requestId string, runtime *runtime.Runtime, template *template.Template,
	functionParam map[string]interface{}) (string, error) {

	var err error
	container := &Container{
		Id:            util.UniqueId(),
		RequestId:     requestId,
		Template:      template.Name,
		FunctionParam: functionParam,
		Timestamp: Timestamp{
			ContainerCreateAt: time.Now().UnixNano() / 1e3,
		},
	}
	service.dataMap[container.Id] = container

	//// 获取 cgroup
	//container.CgroupId, err = service.cgroupService.Get(&template.ResourceLimit)
	//if container.CgroupId == "" {
	//	delete(service.dataMap, container.Id)
	//	return "", errors.New("NewCgroupError")
	//}

	// new file system
	container.BaseFsPath, err = service.fsService.NewContainerFs(container.Id, template.Runtime)
	if err != nil {
		delete(service.dataMap, container.Id)
		service.cgroupService.GiveBack(container.CgroupId)
		return "", errors.New("NewRootFsError")
	}

	// 基于 zygote 创建或者直接启动容器
	if err = service.newContainerProcessByZygote(runtime, template, container); err != nil {
		err = service.newContainerProcessDirectly(runtime, template, container)
	}
	if err != nil {
		delete(service.dataMap, container.Id)
		service.fsService.CleanContainerFs(container.BaseFsPath)
		service.cgroupService.GiveBack(container.CgroupId)
		return "", err
	}

	return container.Id, nil
}

func (service *Service) Get(name string) (container *Container, ok bool) {
	container, ok = service.dataMap[name]
	return
}

func (service *Service) List() map[string]*Container {
	return service.dataMap
}

func (service *Service) Dump() error {
	return util.WriteJsonDataToFile(service.getDataFilePath(), service.dataMap)
}

func (service *Service) getDataFilePath() string {
	return path.Join(config.GetDataPath(), config.ContainerDataFileName)
}
