package service

import (
	"os"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"time"
)

type StoreService interface {
	//存储文件，返回文件标示符
	Save(nameSpace string, content []byte) (string, error)
	//根据文件标示符，获取实际位置
	Get(nameSpace string, id string) []byte
}

//单节点生效，多节点不采用此方式
type LocalStoreService struct {
	//此目录应该只能有读写权限
	Path     string
	duration *time.Duration
}

func (localStoreService *LocalStoreService) Save(nameSpace string, content []byte) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	id := uid.String()
	f, err := os.Create(localStoreService.Path + "/" + nameSpace + "/" + id)
	if f == nil || err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(content)
	return id, err
}
func (localStoreService *LocalStoreService) Get(nameSpace string, id string) []byte {
	f, err := os.Open(localStoreService.Path + "/" + nameSpace + "/" + id)
	if f == nil || err != nil {
		return nil
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	return b
}
