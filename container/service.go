package container

import (
	"engine/config"
	"engine/runtime"
	"engine/template"
	log "github.com/sirupsen/logrus"
	"path"
	"strconv"
	"time"
)

type Service struct {
	fsService     *FsService
	cgroupService *CgroupService
	zygoteService *ZygoteService
}

func NewContainerService(runtimeSet map[string]*runtime.Runtime, templateSet map[string]*template.Template) *Service {

	service := &Service{
		fsService:     NewFsService(),
		cgroupService: NewCgroupService(),
		zygoteService: NewZygoteService(runtimeSet, templateSet),
	}

	return service
}

func (service *Service) Create(id int, runtime *runtime.Runtime, template *template.Template, zygote bool,
	functionParam map[string]interface{}) (string, error) {

	var err error
	container := &Container{
		Id:            strconv.Itoa(id),
		Template:      template.Name,
		FunctionParam: functionParam,
		Timestamp: Timestamp{
			ContainerCreateAt: time.Now().UnixNano() / 1e3,
		},
	}

	// 获取 cgroup
	if template.ResourceLimit != nil {
		container.CgroupId, err = service.cgroupService.Get(template.ResourceLimit)
		if err != nil {
			service.onCreateError(container)
			return "", err
		}
	}

	// new file system
	container.BaseFsPath, err = service.fsService.Get(container.Id, template.Runtime)
	if err != nil {
		service.onCreateError(container)
		return "", err
	}

	// 基于 zygote 创建或者直接启动容器
	if zygote {
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
	if container.BaseFsPath != "" {
		service.fsService.CleanContainerFs(container.BaseFsPath)
	}
	if container.CgroupId != "" {
		service.cgroupService.GiveBack(container.CgroupId)
	}
}

func (service *Service) getDataFilePath() string {
	return path.Join(config.GetDataPath(), config.ContainerDataFileName)
}
