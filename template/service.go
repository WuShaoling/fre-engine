package template

import (
	"engine/config"
	"engine/util"
	"errors"
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

func (service *Service) Create(template *Template) error {
	// 检查 name 是否已存在
	if _, ok := service.dataMap[template.Name]; ok {
		return errors.New("NameExistError")
	}

	// TODO 拉取代码
	// TODO 加载额外所需的共享库和依赖包

	// 更新缓存
	service.dataMap[template.Name] = template

	// 写回文件
	if err := util.WriteJsonDataToFile(service.getDataFilePath(), service.dataMap); err != nil {
		delete(service.dataMap, template.Name)
		return errors.New("SaveDataError")
	}
	return nil
}

func (service *Service) Get(name string) (template *Template, ok bool) {
	template, ok = service.dataMap[name]
	return
}

func (service *Service) List() map[string]*Template {
	return service.dataMap
}

func (service *Service) Delete(name string) error {
	template, ok := service.dataMap[name]
	if !ok {
		return nil
	} else {
		delete(service.dataMap, name)
	}
	if err := util.WriteJsonDataToFile(service.getDataFilePath(), service.dataMap); err != nil {
		service.dataMap[name] = template // 复原
		return err
	}
	return nil
}

func (service *Service) Dump() error {
	return util.WriteJsonDataToFile(service.getDataFilePath(), service.dataMap)
}

func (service *Service) getDataFilePath() string {
	return path.Join(config.GetDataPath(), config.TemplateDataFileName)
}
