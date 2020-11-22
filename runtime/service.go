package runtime

import (
	"engine/config"
	"engine/util"
	log "github.com/sirupsen/logrus"
	"path"
)

type Service struct {
	dataMap map[string]*Runtime
}

func NewRuntimeService() *Service {
	runtime := &Service{
		dataMap: make(map[string]*Runtime),
	}
	util.LoadJsonDataFromFile(runtime.getDataFilePath(), &runtime.dataMap)
	log.Info("start runtime service ok!")
	return runtime
}

func (service *Service) Get(name string) (runtime *Runtime, ok bool) {
	runtime, ok = service.dataMap[name]
	return
}

func (service *Service) List() map[string]*Runtime {
	return service.dataMap
}

func (service *Service) Dump() error {
	return util.WriteJsonDataToFile(service.getDataFilePath(), service.dataMap)
}

func (service *Service) getDataFilePath() string {
	return path.Join(config.GetDataPath(), config.RuntimeDataFileName)
}
