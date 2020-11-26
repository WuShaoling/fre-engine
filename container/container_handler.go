package container

import (
	"engine/config"
	"engine/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

// 容器启动，生成 pid
func (service *Service) ContainerProcessStartHandler(id string, pid int, timestamp int64) {
	//log.Infof("container(id=%s) process(pid=%d) start at %d", id, pid, timestamp)

	// 检查基本信息
	containerInfo := service.getContainerFromStore(id)
	if containerInfo == nil {
		return
	}

	// 修改状态
	containerInfo.Status = config.StatusRunning
	containerInfo.Timestamp.ContainerProcessStartAt = timestamp
	containerInfo.Pid = pid
}

// 容器进程退出，即父进程 wait 方法返回
func (service *Service) ContainerProcessEndHandler(id string, timestamp int64) {
	//log.Infof("container(id=%s) process end at %d", id, timestamp)

	container := service.getContainerFromStore(id)
	if container == nil {
		return
	}

	// 修改状态
	container.Status = config.StatusExit
	container.Timestamp.ContainerProcessEndAt = timestamp

	// 回调容器退出
	service.onContainerExitCallback <- gin.H{
		"id":        id,
		"requestId": container.RequestId,
	}

	// 清理容器环境
	service.fsService.CleanContainerFs(container.BaseFsPath)
	if container.CgroupId != "" {
		service.cgroupService.GiveBack(container.CgroupId)
	}
	container.Timestamp.ContainerDestroyedAt = time.Now().UnixNano() / 1e3

	//fmt.Printf("ContainerCreate --%d--> ContainerProcessStart --%d-->  ContainerProcessRun --%d--> FunctionRun, total=%d\n",
	//	-container.Timestamp.ContainerCreateAt+container.Timestamp.ContainerProcessStartAt,
	//	-container.Timestamp.ContainerProcessStartAt+container.Timestamp.ContainerProcessRunAt,
	//	-container.Timestamp.ContainerProcessRunAt+container.Timestamp.FunctionRunAt,
	//	-container.Timestamp.ContainerCreateAt+container.Timestamp.FunctionRunAt,
	//)
	//fmt.Printf("FunctionRun --%d--> FunctionEnd\n", -container.Timestamp.FunctionRunAt+container.Timestamp.FunctionEndAt)
	//fmt.Printf("FunctionEnd --%d--> ContainerProcessEnd --%d--> ContainerDestroyed, total=%d\n",
	//	-container.Timestamp.FunctionEndAt+container.Timestamp.ContainerProcessEndAt,
	//	-container.Timestamp.ContainerProcessEndAt+container.Timestamp.ContainerDestroyedAt,
	//	-container.Timestamp.FunctionEndAt+container.Timestamp.ContainerDestroyedAt)
	//fmt.Printf("ContainerCreate --%d--> ContainerDestroyed\n", -container.Timestamp.ContainerCreateAt+container.Timestamp.ContainerDestroyedAt)
}

// 函数执行完成，bootstrap 上报结果
func (service *Service) FunctionEndHandler(result model.FunctionEndRequest) {
	//log.Infof("container(id=%s) exec function end", result.Id)

	container := service.getContainerFromStore(result.Id)
	if container == nil {
		return
	}

	container.Timestamp.ContainerProcessRunAt = result.ContainerProcessRunAt
	container.Timestamp.FunctionRunAt = result.FunctionRunTimestamp
	container.Timestamp.FunctionEndAt = result.FunctionEndTimestamp

	functionResult := gin.H{
		"id":        result.Id,
		"requestId": container.RequestId,
		"timestamp": container.Timestamp,
		"record": gin.H{
			"1. prepare container environment(containerCreate->containerProcessStart)": container.Timestamp.ContainerProcessStartAt - container.Timestamp.ContainerCreateAt,
			"2. start bootstrap process(containerProcessStart->containerProcessRun)":   container.Timestamp.ContainerProcessRunAt - container.Timestamp.ContainerProcessStartAt,
			"3. prepare function environment(containerProcessRun->functionRun)":        container.Timestamp.FunctionRunAt - container.Timestamp.ContainerProcessRunAt,
			"4. total(containerCreate->functionRun)":                                   container.Timestamp.FunctionRunAt - container.Timestamp.ContainerCreateAt,
			"5. execute function(functionRun->functionEnd)":                            container.Timestamp.FunctionEndAt - container.Timestamp.FunctionRunAt,
		},
	}
	if result.Error != nil {
		functionResult["error"] = result.Error
	} else {
		functionResult["data"] = result.FunctionResult
	}
	service.onFunctionResultCallback <- functionResult
}

func (service *Service) getContainerFromStore(id string) *Container {
	service.dataMapLock.RLock()
	defer service.dataMapLock.RUnlock()

	container, ok := service.dataMap[id]
	if !ok {
		log.Warnf("container(id=%s) not found", id)
		return nil
	} else if container.Status == config.StatusExit {
		log.Warnf("container(id=%s) already exit", id)
		return nil
	}
	return service.dataMap[id]
}

func (service *Service) setContainerToStore(container *Container) {
	service.dataMapLock.Lock()
	defer service.dataMapLock.Unlock()

	if _, ok := service.dataMap[container.Id]; ok {
		log.Errorf("container(id=%s) exist", container.Id)
		return
	}
	service.dataMap[container.Id] = container
}

func (service *Service) deleteContainerFromStore(id string) {
	service.dataMapLock.Lock()
	defer service.dataMapLock.Unlock()
	delete(service.dataMap, id)
}
