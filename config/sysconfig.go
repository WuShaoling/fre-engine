package config

import (
	"engine/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

type SysConfig struct {
	RootPath             string `yaml:"rootPath"`             // 应用根目录
	ZygoteMaxMemory      int    `yaml:"zygoteMaxMemory"`      // zygote 池最大内存，MB
	CgroupPoolSize       int    `yaml:"cgroupPoolSize"`       // cgroup 缓存池大小
	RootfsPoolSize       int    `yaml:"rootfsPoolSize"`       // rootfs 缓存池大小
	ServePort            string `yaml:"servePort"`            // 服务监听的地址
	ContainerCodePath    string `yaml:"containerCodePath"`    // 容器内代码的根目录
	ZygoteUnixSocketFile string `yaml:"zygoteUnixSocketFile"` // 容器内代码的根目录
	EnableZygote         bool   `yaml:"enableZygote"`         // 是否开启 Zygote
}

var SysConfigInstance *SysConfig

func InitSysConfig(configPath string) {
	if configPath == "" {
		useDefaultConfig()
	} else {
		useConfigFromFile(configPath)
	}
	initFilePath()
	log.Infof("config: %+v", *SysConfigInstance)
}

func useDefaultConfig() {
	SysConfigInstance = &SysConfig{
		RootPath:             "./workspace",
		ZygoteMaxMemory:      1024,
		CgroupPoolSize:       256,
		RootfsPoolSize:       256,
		ServePort:            "80",
		ContainerCodePath:    "/code",
		ZygoteUnixSocketFile: "/tmp/free.zygote.sock",
		EnableZygote:         true,
	}
}

func useConfigFromFile(configPath string) {
	SysConfigInstance = &SysConfig{}

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatal("init error, ", err)
	}

	if err := yaml.NewDecoder(f).Decode(SysConfigInstance); err != nil {
		log.Fatal("load yaml config error, ", err)
	}
}

func initFilePath() {
	_ = os.Remove(SysConfigInstance.ZygoteUnixSocketFile)

	if SysConfigInstance.RootPath[0] != '/' {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal("initFilePath error, ", err)
		}
		SysConfigInstance.RootPath = path.Join(pwd, SysConfigInstance.RootPath)
	}

	if err := util.MkdirIfNotExist(GetLogPath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetDataPath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetRuntimePath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetVolumeHostPath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetZygoteCodePath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetContainerFsPath()); err != nil {
		log.Fatal("[error]", err)
	}
}

func GetLogPath() string         { return path.Join(SysConfigInstance.RootPath, LogPath) }
func GetDataPath() string        { return path.Join(SysConfigInstance.RootPath, DataPath) }
func GetRuntimePath() string     { return path.Join(SysConfigInstance.RootPath, RuntimePath) }
func GetVolumeHostPath() string  { return path.Join(SysConfigInstance.RootPath, VolumeHostPath) }
func GetZygoteCodePath() string  { return path.Join(SysConfigInstance.RootPath, ZygoteCodePath) }
func GetContainerFsPath() string { return path.Join(SysConfigInstance.RootPath, ContainerFsPath) }
