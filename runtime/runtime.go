package runtime

type Runtime struct { // 基础环境
	Name          string   `json:"name" binding:"required"`
	Entrypoint    []string `json:"entrypoint" binding:"required"`
	ZygoteCommand []string `json:"zygoteCommand"`
}
