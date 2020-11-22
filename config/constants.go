package config

const (
	LogPath         = "log"       // 日志目录
	DataPath        = "metadata"  // 数据目录
	RuntimePath     = "runtime"   // 运行时环境目录
	VolumeHostPath  = "volume"    // 数据卷主机端目录
	ContainerFsPath = "container" // 容器文件系统目录
	ZygoteCodePath  = "zygote"    // zygote代码目录
)

const (
	RuntimeDataFileName   = "runtime.json"
	TemplateDataFileName  = "template.json"
	ContainerDataFileName = "container.json"
)

const (
	StatusRunning = "running"
	StatusExit    = "exit"
)
