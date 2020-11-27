package container

type Container struct { // 容器实例
	Id            string                 `json:"id"`         // 唯一标识
	RequestId     string                 `json:"requestId"`  // 请求标识
	Template      string                 `json:"template"`   // 函数模版
	Pid           int                    `json:"pid"`        // 进程ID
	CgroupId      string                 `json:"cgroupId"`   // cgroupId
	Status        string                 `json:"status"`     // 状态
	BaseFsPath    string                 `json:"baseFsPath"` // 文件系统目录 containerBasePath/id
	FunctionParam map[string]interface{} `json:"Param"`      // 函数参数
	Timestamp     Timestamp              `json:"timestamp"`  // 各个阶段的时间戳
}

type Timestamp struct {
	ContainerCreateAt       int64 `json:"containerCreateAt"`       // 开始准备容器环境的时间
	ContainerProcessStartAt int64 `json:"containerProcessStartAt"` // 容器进程创建成功的时间（生成了pid）
	ContainerProcessRunAt   int64 `json:"containerRunAt"`          // 容器进程开始运行的时间
	FunctionRunAt           int64 `json:"functionRunAt"`           // 函数开始运行的时间
	FunctionEndAt           int64 `json:"functionEndAt"`           // 函数运行结束的时间
	ContainerProcessEndAt   int64 `json:"containerEndAt"`          // 容器进程运行结束的时间，即父进程中wait返回时的时间戳
	ContainerDestroyedAt    int64 `json:"containerDestroyedAt"`    // 容器销毁的时间
}

// 运行函数的上下文
type FunctionExecContext struct {
	Id                string                 `json:"id"`         // 函数实例 id
	CodePath          string                 `json:"codePath"`   // 函数代码的路径，/code/templateName
	Handler           string                 `json:"handler"`    // 代码的入口文件
	Params            map[string]interface{} `json:"params"`     // 函数的参数
	ServePort         string                 `json:"servePort"`  // server 的地址，用于和server通信
	RootFsPath        string                 `json:"rootFsPath"` // 根文件系统路径
	CgroupId          string                 `json:"cgroupId"`
	ContainerCreateAt int64                  `json:"containerCreateAt"`
}
