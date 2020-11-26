package container

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Exec() {
	// read command from pipe
	command, err := readFunctionContextFromPipe()
	if err != nil {
		return
	}

	// chroot
	if err := setUpMount(); err != nil {
		return
	}

	// 查找可执行文件
	path, err := exec.LookPath(command[0])
	if err != nil {
		log.Error("exec loop path error: ", err)
		return
	}

	if err := syscall.Exec(path, command[0:], os.Environ()); err != nil {
		log.Error("exec error: ", err)
		return
	}
}

func readFunctionContextFromPipe() ([]string, error) {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()

	data, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Error("exec read pipe error: ", err)
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
		log.Error("get current location error: ", err)
		return err
	}

	err = syscall.Chroot(pwd)
	if err != nil {
		log.Errorf("chroot to %s error %+v", pwd, err)
	}
	return err
}
