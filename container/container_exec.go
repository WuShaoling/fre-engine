package container

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Exec() error {
	// read command from pipe
	command, err := readFunctionContextFromPipe()
	if err != nil {
		return err
	}

	// chroot
	if err := setUpMount(); err != nil {
		return err
	}

	// 查找可执行文件
	path, err := exec.LookPath(command[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}

	if err := syscall.Exec(path, command[0:], os.Environ()); err != nil {
		log.Errorf("syscall.Exec failed, error=%v", err.Error())
		return err
	}
	return nil
}

func readFunctionContextFromPipe() ([]string, error) {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()

	data, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("exec read pipe error %v", err)
		return nil, err
	}

	command := string(data)
	separatorIndex := strings.Index(command, "|")
	entrypoint := command[:separatorIndex]
	entrypointParam := command[separatorIndex+1:]

	// ["python3", "bootstrap.py", "json化后的参数，可以包含空格"]
	return append(strings.Split(entrypoint, " "), entrypointParam), nil
}

func setUpMount() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return err
	}

	err = syscall.Chroot(pwd)
	if err != nil {
		log.Errorf("chroot to %s error %v", pwd, err)
	}
	return err

	//merge proc
	//defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	//syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	//syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}
