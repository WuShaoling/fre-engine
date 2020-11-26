package container

import (
	"engine/config"
	"engine/runtime"
	"engine/template"
	"engine/util"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"path"
	"sync"
	"time"
)

type Service struct {
	onContainerExitCallback  chan gin.H
	onFunctionResultCallback chan gin.H

	dataMap     map[string]*Container
	dataMapLock sync.RWMutex

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

	log.Info("start container service ok!")
	return service
}

func (service *Service) Create(requestId string, runtime *runtime.Runtime, template *template.Template,
	zygote string, functionParam map[string]interface{}) (string, error) {

	log.Infof("create container: requestId=%s, runtime=%s, template=%s", requestId, runtime.Name, template.Name)

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
	service.setContainerToStore(container)

	// 获取 cgroup
	if template.ResourceLimit != nil {
		container.CgroupId, err = service.cgroupService.Get(template.ResourceLimit)
		if err != nil {
			service.onCreateError(container)
			return "", err
		}
	}

	// new file system
	t1 := time.Now().UnixNano() / 1e3
	container.BaseFsPath = service.fsService.Get(container.Id, template.Runtime)
	t2 := time.Now().UnixNano() / 1e3
	fmt.Println("---->service.fsService.Get(): ", t2-t1)
	//// new file system
	//container.BaseFsPath, err = service.fsService.Get(container.Id, template.Runtime)
	//if err != nil {
	//	service.onCreateError(container)
	//	return "", err
	//}

	// 基于 zygote 创建或者直接启动容器
	if zygote == "true" {
		if err = service.newContainerProcessByZygote(runtime, template, container); err != nil {
			log.Warnf("new container by zygote failed, error=%+v", err)
			err = service.newContainerProcessDirectly(runtime, template, container)
		}
	} else {
		err = service.newContainerProcessDirectly(runtime, template, container)
	}

	if err != nil {
		service.onCreateError(container)
		return "", err
	}

	return container.Id, nil
}

func (service *Service) onCreateError(container *Container) {
	service.deleteContainerFromStore(container.Id)
	if container.BaseFsPath != "" {
		service.fsService.CleanContainerFs(container.BaseFsPath)
	}
	if container.CgroupId != "" {
		service.cgroupService.GiveBack(container.CgroupId)
	}
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
