package container

import (
	"encoding/json"
	"engine/config"
	"engine/runtime"
	"engine/template"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// 直接启动容器进程

func (service *Service) newContainerProcessDirectly(runtime *runtime.Runtime, template *template.Template, container *Container) error {
	//log.Infof("new container process(%s) directly, runtime=%s, template=%s", container.Id, runtime.Name, template.Name)

	// new pipe
	//t1 := time.Now().UnixNano() / 1e3
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorf("new container(id=%s) process directly: new pipe error, %+v", container.Id, err)
		return err
	}
	defer writePipe.Close()
	//t2 := time.Now().UnixNano() / 1e3
	//fmt.Println("---->os.Pipe: ", t2-t1)

	// new process
	//t1 = time.Now().UnixNano() / 1e3
	containerProcess := exec.Command("/proc/self/exe", "exec")
	containerProcess.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
	}
	containerProcess.ExtraFiles = []*os.File{readPipe}
	containerProcess.Dir = service.fsService.GetContainerRootFsPath(container.BaseFsPath)
	containerProcess.Stdout = os.Stdout
	containerProcess.Stderr = os.Stderr
	containerProcess.Env = append(
		os.Environ(),
		fmt.Sprintf("id=%s", container.Id),
		fmt.Sprintf("containerCreateAt=%d", container.Timestamp.ContainerCreateAt),
		fmt.Sprintf("containerProcessStartAt=%d", time.Now().UnixNano()/1e3),
	)
	//t2 = time.Now().UnixNano() / 1e3
	//fmt.Println("---->exec.Command: ", t2-t1)

	// start process
	//t1 = time.Now().UnixNano() / 1e3
	if err := containerProcess.Start(); err != nil {
		log.Errorf("new container(id=%s) process directly: start process error, %+v", container.Id, err)
		service.onNewContainerProcessDirectlyError(readPipe, nil)
		return err
	}
	//t2 = time.Now().UnixNano() / 1e3
	//fmt.Println("---->containerProcess.Start(): ", t2-t1)

	// 加入 cgroup
	if container.CgroupId != "" {
		if err := service.cgroupService.Set(container.CgroupId, containerProcess.Process.Pid); err != nil {
			log.Errorf("set container(%s) to cgroup(id=%s) error, %+v", container.Id, container.CgroupId, err)
			service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
			return err
		}
	}

	// 构建函数执行上下文
	//t1 = time.Now().UnixNano() / 1e3
	functionExecContext := service.buildFunctionExecContext(template, container)
	if functionExecContext == "" {
		service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
		return errors.New("GetFunctionExecContextError")
	}

	// 发送运行命令，例如：python3 bootstrap.py|param str, |为解析命令的分隔符
	command := strings.Join(runtime.Entrypoint, " ") + "|" + functionExecContext
	if _, err := writePipe.Write([]byte(command)); err != nil {
		log.Errorf("new container(id=%s) process directly: send init command error, %+v", command, err)
		service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
		return err
	}
	//t2 = time.Now().UnixNano() / 1e3
	//fmt.Println("---->writePipe.Write(): ", t2-t1)

	// 异步 wait 容器进程退出
	go func() {
		_ = containerProcess.Wait()
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
		Id:                container.Id,
		CodePath:          service.fsService.GetContainerFunctionCodePath(template.Name),
		Handler:           template.Handler,
		Params:            container.FunctionParam,
		ServePort:         config.SysConfigInstance.ServePort,
		RootFsPath:        service.fsService.GetContainerRootFsPath(container.BaseFsPath),
		CgroupId:          container.CgroupId,
		ContainerCreateAt: container.Timestamp.ContainerCreateAt,
	}

	data, err := json.Marshal(ctx)
	if err != nil {
		log.Error("build function exec context: json marshal error, ", err)
		return ""
	}
	return string(data)
}
