package template

import (
	"engine/config"
	"engine/util"
	log "github.com/sirupsen/logrus"
	"path"
)

type Service struct {
	dataMap map[string]*Template
}

func NewTemplateService() *Service {
	service := &Service{
		dataMap: make(map[string]*Template),
	}
	util.LoadJsonDataFromFile(service.getDataFilePath(), &service.dataMap)
	log.Info("start template service ok!")
	return service
}

func (service *Service) Get(name string) *Template {
	return service.dataMap[name]
}

func (service *Service) List() map[string]*Template {
	return service.dataMap
}

func (service *Service) getDataFilePath() string {
	return path.Join(config.GetDataPath(), config.TemplateDataFileName)
}
