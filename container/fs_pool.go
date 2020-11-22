package container

import (
	"engine/config"
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
}

func NewFsService() *FsService {
	log.Info("start fs service ok!")
	return &FsService{
	}
}

func (service *FsService) NewContainerFs(id, runtime string) (string, error) {
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
