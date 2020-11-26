package container

import (
	"engine/config"
	"engine/util"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"syscall"
)

const (
	MergePath  = "merge"
	UpperPath  = "upper"
	WorkerPath = "worker"
)

type FsService struct {
	fsPool chan string
}

func NewFsService() *FsService {
	log.Info("start fs service ok!")
	service := &FsService{
		fsPool: make(chan string, config.SysConfigInstance.RootfsPoolSize),
	}
	service.initPool()
	return service
}

func (service *FsService) Get(id, runtime string) string {
	return <-service.fsPool
}

func (service *FsService) newContainerFs(id, runtime string) (string, error) {
	basePath := path.Join(config.GetContainerFsPath(), id)
	lowerPath := path.Join(config.GetRuntimePath(), runtime)
	mergePath := path.Join(basePath, MergePath)
	upperPath := path.Join(basePath, UpperPath)
	workerPath := path.Join(basePath, WorkerPath)

	// 创建目录
	if err := os.Mkdir(basePath, 0755); err != nil {
		log.Errorf("new basePath %s for container %s error, %+v", basePath, id, err)
		_ = os.RemoveAll(basePath)
		return "", err
	}
	if err := os.Mkdir(mergePath, 0755); err != nil {
		log.Errorf("new mergePath %s for container %s error, %+v", mergePath, id, err)
		_ = os.RemoveAll(basePath)
		return "", err
	}
	if err := os.Mkdir(upperPath, 0755); err != nil {
		log.Errorf("new upperPath %s for container %s error, %+v", upperPath, id, err)
		_ = os.RemoveAll(basePath)
		return "", err
	}
	if err := os.Mkdir(workerPath, 0755); err != nil {
		log.Errorf("new workerPath %s for container %s error, %+v", workerPath, id, err)
		_ = os.RemoveAll(basePath)
		return "", err
	}

	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerPath, upperPath, workerPath)
	if err := syscall.Mount("overlay", mergePath, "overlay", 0, data); err != nil {
		log.Errorf("overlay mount for container %s error, data=%s, mountPath=%s, error=%v", id, data, mergePath, err)
		_ = os.RemoveAll(basePath)
		return "", err
	}

	// TODO 挂载数据目录
	return basePath, nil
}

func (service *FsService) CleanContainerFs(basePath string) {
	_ = syscall.Unmount(path.Join(basePath, MergePath), 0)
	_ = os.RemoveAll(basePath)
}

func (service *FsService) GetContainerRootFsPath(basePath string) string {
	return path.Join(basePath, MergePath)
}

func (service *FsService) GetContainerFunctionCodePath(templateName string) string {
	return path.Join(config.SysConfigInstance.ContainerCodePath, templateName)
}

func (service *FsService) initPool() {
	for i := 0; i < config.SysConfigInstance.RootfsPoolSize; i++ {
		basePath, err := service.newContainerFs(util.UniqueId(), "python3.7")
		if err != nil {
			log.Fatal(err)
		} else {
			service.fsPool <- basePath
		}
	}

	for i := 0; i < 4; i++ { // 4个生产者同时生产
		go func(id int) {
			errorCount := 0
			for {
				basePath, err := service.newContainerFs(util.UniqueId(), "python3.7")
				if err == nil {
					service.fsPool <- basePath
					errorCount = 0
				} else {
					errorCount++
					if errorCount > 16 {
						log.Errorf("fs producer(id=%d) too many errors, exit", id)
						return
					}
				}
			}
		}(i)
	}
}
