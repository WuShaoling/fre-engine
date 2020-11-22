package container

import (
	"encoding/json"
	"engine/config"
	"engine/runtime"
	"engine/template"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// 直接启动容器进程

func (service *Service) newContainerProcessDirectly(runtime *runtime.Runtime, template *template.Template, container *Container) error {
	// new pipe
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorf("newContainerProcessDirectly(id=%s) new pipe error, %+v", container.Id, err)
		return err
	}
	defer writePipe.Close()

	// new process
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Errorf("newContainerProcessDirectly get init process error %v", err)
		return err
	}
	containerProcess := exec.Command(initCmd, "exec")
	containerProcess.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
	}
	containerProcess.ExtraFiles = []*os.File{readPipe}
	containerProcess.Env = append(os.Environ(), template.Envs...)
	containerProcess.Dir = service.fsService.GetContainerRootFsPath(container.BaseFsPath)
	// 临时使用std, TODO 重定向 Stdout Stderr 到日志文件
	containerProcess.Stdout = os.Stdout
	containerProcess.Stderr = os.Stderr
	if err := containerProcess.Start(); err != nil {
		log.Errorf("newContainerProcessDirectly(id=%s) start process error, %+v", container.Id, err)
		service.onNewContainerProcessDirectlyError(readPipe, nil)
		return err
	}

	// 容器进程创建成功
	service.ContainerProcessStartHandler(container.Id, containerProcess.Process.Pid, time.Now().UnixNano()/1e3)

	//// 加入 cgroup
	//if err := service.cgroupService.Set(container.CgroupId, containerProcess.Process.Pid); err != nil {
	//	service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
	//	return err
	//}

	// 构建函数执行上下文
	functionExecContext := service.buildFunctionExecContext(template, container)
	if functionExecContext == "" {
		service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
		return errors.New("GetFunctionExecContextError")
	}

	// 发送运行命令，例如：python3 bootstrap.py|param str, |为解析命令的分隔符
	command := strings.Join(runtime.Entrypoint, " ") + "|" + functionExecContext
	if _, err := writePipe.Write([]byte(command)); err != nil {
		log.Errorf("sendInitCommand failed, command=%s, error=%v", command, err)
		service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
		return err
	}

	// 异步 wait 容器进程退出
	go func() {
		if err := containerProcess.Wait(); err != nil {
			log.Errorf("wait error, id=%s, pid=%d, error=%+v", container.Id, containerProcess.Process.Pid, err)
		}
		service.ContainerProcessEndHandler(container.Id, time.Now().UnixNano()/1e3)
	}()

	return nil
}

func (service *Service) onNewContainerProcessDirectlyError(readPipe *os.File, containerProcess *exec.Cmd) {
	if readPipe != nil {
		_ = readPipe.Close()
	}

	if containerProcess != nil && containerProcess.Process != nil {
		_ = containerProcess.Process.Kill()
		_ = containerProcess.Wait()
	}
}

func (service *Service) buildFunctionExecContext(template *template.Template, container *Container) string {
	ctx := FunctionExecContext{
		Id:         container.Id,
		CodePath:   service.fsService.GetContainerFunctionCodePath(template.Name),
		Handler:    template.Handler,
		Params:     container.FunctionParam,
		ServePort:  config.SysConfigInstance.ServePort,
		RootFsPath: service.fsService.GetContainerRootFsPath(container.BaseFsPath),
		CgroupId:   container.CgroupId,
	}

	data, err := json.Marshal(ctx)
	if err != nil {
		log.Error("buildFunctionExecContext json marshal error, ", err)
		return ""
	}
	return string(data)
}
