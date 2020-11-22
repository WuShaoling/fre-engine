package container

import (
	"engine/config"
	"engine/model"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

// 容器启动，生成 pid
func (service *Service) ContainerProcessStartHandler(id string, pid int, timestamp int64) {
	log.Infof("ContainerProcessStartHandler: id=%s, pid=%d, timestamp=%d", id, pid, timestamp)

	// 检查基本信息
	containerInfo := service.getAndCheckContainer(id)
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
	log.Infof("ContainerProcessEndHandler: id=%s, timestamp=%d", id, timestamp)

	container := service.getAndCheckContainer(id)
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

	fmt.Printf("ContainerCreate --%d--> ContainerProcessStart --%d-->  ContainerProcessRun --%d--> FunctionRun, total=%d\n",
		-container.Timestamp.ContainerCreateAt+container.Timestamp.ContainerProcessStartAt,
		-container.Timestamp.ContainerProcessStartAt+container.Timestamp.ContainerProcessRunAt,
		-container.Timestamp.ContainerProcessRunAt+container.Timestamp.FunctionRunAt,
		-container.Timestamp.ContainerCreateAt+container.Timestamp.FunctionRunAt,
	)
	fmt.Printf("FunctionRun --%d--> FunctionEnd\n", -container.Timestamp.FunctionRunAt+container.Timestamp.FunctionEndAt)
	fmt.Printf("FunctionEnd --%d--> ContainerProcessEnd --%d--> ContainerDestroyed, total=%d\n",
		-container.Timestamp.FunctionEndAt+container.Timestamp.ContainerProcessEndAt,
		-container.Timestamp.ContainerProcessEndAt+container.Timestamp.ContainerDestroyedAt,
		-container.Timestamp.FunctionEndAt+container.Timestamp.ContainerDestroyedAt)
	fmt.Printf("ContainerCreate --%d--> ContainerDestroyed\n", -container.Timestamp.ContainerCreateAt+container.Timestamp.ContainerDestroyedAt)
}

// 函数执行完成，bootstrap 上报结果
func (service *Service) FunctionEndHandler(result model.FunctionEndRequest) {
	log.Infof("FunctionEndHandler: id=%s", result.Id)

	container := service.getAndCheckContainer(result.Id)
	if container == nil {
		return
	}

	container.Timestamp.ContainerProcessRunAt = result.ContainerProcessRunAt
	container.Timestamp.FunctionRunAt = result.FunctionRunTimestamp
	container.Timestamp.FunctionEndAt = result.FunctionEndTimestamp

	functionResult := gin.H{
		"id":        result.Id,
		"requestId": container.RequestId,
	}
	if result.Error != nil {
		functionResult["error"] = result.Error
	} else {
		functionResult["data"] = result.FunctionResult
	}
	service.onFunctionResultCallback <- functionResult
}

func (service *Service) getAndCheckContainer(id string) *Container {
	container, ok := service.dataMap[id]
	if !ok {
		log.Warnf("container(id=%s) not found", id)
		return nil
	} else if container.Status == config.StatusExit {
		log.Warnf("container(id=%s) already exit", id)
		return nil
	}
	return container
}
