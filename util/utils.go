package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func MkdirIfNotExist(p string) error {
	if f, err := os.Stat(p); err == nil { // 如果已经存在了
		if !f.IsDir() { // 如果不是目录，报错
			return errors.New(fmt.Sprintf("%s exist and is a file", p))
		}
	} else if os.IsNotExist(err) { // 如果不存在
		if err := os.MkdirAll(p, 0744); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func WriteToFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		//log.Errorf("WriteToFile open file %s error, %v", filename, err)
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		//log.Errorf("WriteToFile write to file %s error, %v", filename, err)
	}
	return err
}

func ReadFromFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		//log.Errorf("ReadFromFile open file %s error, %v", filename, err)
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		//log.Errorf("ReadFromFile read file %s error, %v", filename, err)
	}
	return data, err
}

func WriteJsonDataToFile(filename string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		log.Errorf("WriteJsonDataToFile json.Marshal error, filename=%s, data=%v, error=%v", filename, v, err)
		return err
	}
	return WriteToFile(filename, data)
}

func LoadJsonDataFromFile(filename string, v interface{}) {
	if data, err := ReadFromFile(filename); err == nil {
		if e := json.Unmarshal(data, v); e != nil {
			log.Errorf("LoadJsonDataFromFile json.Unmarshal error, filename=%s, error=%v", filename, err)
		}
	}
	// 文件可能不存在，直接忽略
}

func UniqueId() string {
	return uuid.New().String()
	//return strconv.FormatInt(time.Now().UnixNano()/1e3, 10)
}

func Int16ToBytes(n int16) ([]byte, error) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytesBuffer, binary.LittleEndian, n); err != nil {
		log.Error("IntToBytes error, ", err)
		return nil, err
	}
	return bytesBuffer.Bytes(), nil
}
