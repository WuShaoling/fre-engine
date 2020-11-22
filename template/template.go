package template

type Template struct { // 函数模版
	Name    string `json:"name" binding:"required"`    // 函数名, 唯一
	Version string `json:"version" binding:"required"` // 版本

	Runtime    string   `json:"runtime" binding:"required"` // 基础环境
	Handler    string   `json:"handler" binding:"required"` // 函数入口文件
	Packages   []string `json:"packages"`                   // 依赖包
	SharedLibs []string `json:"sharedLibs"`                 // 基础环境之外所需的共享库
	Volume     string   `json:"volume"`                     // 数据卷挂载目录
	Envs       []string `json:"envs"`                       // 环境变量, key=value

	ResourceLimit ResourceLimit `json:"resourceLimit"` // 资源限制
}

type ResourceLimit struct {
	Memory   int64  `json:"memory"`
	CpuShare uint64 `json:"cpuShares"`
	Timeout  int    `json:"timeout"` // 超时时间
}
