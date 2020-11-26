package container

import (
	"engine/config"
	"engine/template"
	"engine/util"
	"errors"
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
	"sync"
)

const cgroupPrefix = "/fre_"

type CgroupService struct {
	pool    []string
	dataMap map[string]*cgroups.Cgroup
	mutex   sync.Mutex
}

func NewCgroupService() *CgroupService {
	service := &CgroupService{
		pool:    make([]string, 0, config.SysConfigInstance.CgroupPoolSize),
		dataMap: make(map[string]*cgroups.Cgroup),
	}

	// 构造缓存池
	for i := 0; i < config.SysConfigInstance.CgroupPoolSize; i++ {
		if id, err := service.newCgroup(nil); id == "" {
			log.Fatal("NewCgroupService: new cgroup error", err)
		} else {
			service.pool = append(service.pool, id)
		}
	}

	log.Info("start cgroup service ok!")
	return service
}

// 获取 cgroup
func (service *CgroupService) Get(limit *template.ResourceLimit) (string, error) {
	service.mutex.Lock()

	n := len(service.pool)

	if n > 0 {
		c := service.pool[n-1]
		service.pool = service.pool[0 : n-1]
		service.mutex.Unlock()
		return c, nil
	}

	service.mutex.Unlock()
	return service.newCgroup(limit)
}

func (service *CgroupService) GiveBack(id string) {
	cgroup := service.getOrLoad(id)
	if cgroup != nil {
		if len(service.pool) < config.SysConfigInstance.CgroupPoolSize {
			service.mutex.Lock()
			service.pool = append(service.pool, id)
			service.mutex.Unlock()
		} else {
			_ = (*cgroup).Delete()
		}
	}
}

func (service *CgroupService) Set(id string, pid int) error {
	c := service.getOrLoad(id)
	if c == nil {
		return errors.New(fmt.Sprintf("cgroup(id=%s) not found", id))
	}

	err := (*c).Add(cgroups.Process{Pid: pid})
	if err != nil {
		log.Errorf("add pid(%d) to cgroup(%s) error, %+v", pid, id, err)
	}
	return err
}

func (service *CgroupService) getOrLoad(id string) *cgroups.Cgroup {
	// 先从map找
	if c, ok := service.dataMap[id]; ok {
		return c
	}

	// map 找不到从本地 load
	c, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(cgroupPrefix+id))
	if err != nil {
		log.Errorf("load cgroup(id=%s) error, %+v", id, err)
		return nil
	}

	service.dataMap[id] = &c
	return &c
}

func (service *CgroupService) newCgroup(limit *template.ResourceLimit) (string, error) {
	id := util.UniqueId()

	linuxResource := &specs.LinuxResources{
		Memory: &specs.LinuxMemory{},
		CPU:    &specs.LinuxCPU{},
	}
	if limit != nil {
		linuxResource.Memory.Limit = &limit.Memory
		linuxResource.CPU.Shares = &limit.CpuShare
	}

	cgroup, err := cgroups.New(cgroups.V1, cgroups.StaticPath(cgroupPrefix+id), linuxResource)
	if err != nil {
		log.Error("new cgroup error: ", err)
		return "", err
	}

	service.dataMap[id] = &cgroup
	return id, nil
}
