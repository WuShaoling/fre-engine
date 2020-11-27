package container

import (
	"bufio"
	"encoding/json"
	"engine/config"
	"engine/runtime"
	"engine/template"
	"engine/util"
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/exec"
	"path"
)

type ZygoteProcess struct {
	id         string
	cmd        *exec.Cmd
	packageSet []string
}

type ZygoteProcessTree map[string]*ZygoteProcess // key 暂时用 templateName，实际应该为 ZygoteProcess.id

type ZygoteService struct {
	zygoteProcessUnixSocket  map[string]*net.UnixConn     // key 为 id
	runtimeZygoteProcessTree map[string]ZygoteProcessTree // key 为 runtimeName
}

func NewZygoteService(runtimeSet map[string]*runtime.Runtime, templateSet map[string]*template.Template) *ZygoteService {
	service := &ZygoteService{
		zygoteProcessUnixSocket:  make(map[string]*net.UnixConn),
		runtimeZygoteProcessTree: make(map[string]ZygoteProcessTree),
	}

	if err := service.startUnixSocketServer(); err != nil {
		log.Fatal("NewZygoteService startUnixSocketServer error, ", err)
	}
	service.buildRuntimeZygoteProcessTree(runtimeSet, templateSet)

	log.Info("start zygote service ok!")
	return service
}

func (service *ZygoteService) NewContainerByZygoteProcess(runtimeName, templateName string, messageBody string) error {
	log.Infof("new container by zygote process, runtime=%s, template=%s", runtimeName, templateName)

	runtimeZygoteProcessTree, ok := service.runtimeZygoteProcessTree[runtimeName]
	if !ok {
		return errors.New("RuntimeNotSupportZygote")
	}

	zygoteProcess, ok := runtimeZygoteProcessTree[templateName]
	if !ok {
		return errors.New("NoMatchZygoteProcessFound")
	}

	//log.Infof("find matched zygote process, id=%s, pid=%d", zygoteProcess.id, zygoteProcess.cmd.Process.Pid)
	unixSocket, ok := service.zygoteProcessUnixSocket[zygoteProcess.id]
	if !ok {
		return errors.New("ZygoteProcessUnixSocketNotFound")
	}

	msgHeader, err := util.Int16ToBytes(int16(len(messageBody)))
	if err != nil {
		return err
	}
	_, err = unixSocket.Write(append(msgHeader, []byte(messageBody)...))
	if err != nil {
		log.Errorf("send command to zygoteProcess(id=%s, pid=%d) error, %+v",
			zygoteProcess.id, zygoteProcess.cmd.Process.Pid, err)
	}
	//log.Infof("send command to zygoteProcess(id=%s, pid=%d) ok", zygoteProcess.id, zygoteProcess.cmd.Process.Pid)
	return err
}

func (service *ZygoteService) startUnixSocketServer() error {
	unixAddr, err := net.ResolveUnixAddr("unix", config.SysConfigInstance.ZygoteUnixSocketFile)
	if err != nil {
		log.Error("zygote service: resolve unix address error, ", err)
		return err
	}

	unixListener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		log.Error("zygote service: listen unix error, ", err)
		return err
	}

	go func() {
		for {
			unixConn, err := unixListener.AcceptUnix()
			if err != nil {
				log.Error("zygote service: accept unix error, ", err)
				continue
			}
			// 异步等待 zygote process 注册
			go func() {
				reader := bufio.NewReader(unixConn)
				if id, err := reader.ReadString('\n'); err != nil {
					log.Error("zygote service: read message error, ", err)
					return
				} else {
					id = id[:len(id)-1]
					log.Infof("zygote service: receive zygote process(id=%s) register", id)
					service.zygoteProcessUnixSocket[id] = unixConn
				}
			}()
		}
	}()
	return nil
}

// 创建 runtime 的 zygoteProcess 树
func (service *ZygoteService) buildRuntimeZygoteProcessTree(runtimeSet map[string]*runtime.Runtime, templateSet map[string]*template.Template) {
	// 对 templateSet 按 runtime 分组
	templateGroup := make(map[string][]*template.Template)
	for _, v := range templateSet {
		templateGroup[v.Runtime] = append(templateGroup[v.Runtime], v)
	}

	// 对于每一种 runtime，构造 ZygoteProcessTree
	for runtimeName, templateList := range templateGroup {
		r, ok := runtimeSet[runtimeName]
		if !ok || r.ZygoteCommand == nil && len(r.ZygoteCommand) == 0 {
			continue
		}
		zygoteProcessTree := ZygoteProcessTree{}
		for _, t := range templateList {
			if zygoteProcess := service.newZygoteProcess(r, t); zygoteProcess != nil {
				zygoteProcessTree[zygoteProcess.id] = zygoteProcess
				go service.onZygoteProcessExit(r.Name, zygoteProcess) // 异步等待进程退出
			}
		}
		service.runtimeZygoteProcessTree[runtimeName] = zygoteProcessTree
	}
}

// 创建 zygote 进程
func (service *ZygoteService) newZygoteProcess(r *runtime.Runtime, t *template.Template) *ZygoteProcess {
	id := t.Name

	// 构造进程参数
	param := map[string]interface{}{
		"id":               id,
		"packageSet":       t.Packages,
		"serverSocketFile": config.SysConfigInstance.ZygoteUnixSocketFile,
	}
	data, err := json.Marshal(param)
	if err != nil {
		log.Errorf("new zygote process: build param error, runtime=%s, template=%s, error=%+v", r.Name, t.Name, err)
		return nil
	}

	// 启动进程
	cmd := exec.Command(r.ZygoteCommand[0], append(r.ZygoteCommand[1:], string(data))...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = path.Join(config.GetZygoteCodePath(), r.Name)
	if err := cmd.Start(); err != nil {
		log.Errorf("new zygote process: start process error, runtime=%s, template=%s, error=%+v", r.Name, t.Name, err)
		return nil
	}

	log.Infof("new zygote process: new zygote process for runtime(%s) template(%s), pid=%d", r.Name, t.Name, cmd.Process.Pid)
	return &ZygoteProcess{
		id:         id,
		cmd:        cmd,
		packageSet: t.Packages,
	}
}

// 进程退出了
func (service *ZygoteService) onZygoteProcessExit(runtimeName string, zygoteProcess *ZygoteProcess) {
	_ = zygoteProcess.cmd.Wait()
	log.Infof("on zygote process(id=%s) exit, runtime=%s", zygoteProcess.id, runtimeName)

	if unixSock, ok := service.zygoteProcessUnixSocket[zygoteProcess.id]; ok {
		_ = unixSock.Close()
		delete(service.zygoteProcessUnixSocket, zygoteProcess.id)
	}

	if zygoteProcessTree, ok := service.runtimeZygoteProcessTree[runtimeName]; ok {
		delete(zygoteProcessTree, zygoteProcess.id)
	}
	// TODO 触发树的调整
}
