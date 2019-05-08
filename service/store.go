package service

import (
	"os"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"heart/entity"
	"time"
	"github.com/astaxie/beego/logs"
)

type StoreService interface {
	//存储文件，返回文件标示符
	Save(nameSpace string, content *[]byte, suffix string) (string, error)
	//根据文件标示符，获取实际位置
	Get(nameSpace string, url string) ([]byte, string, error)
	//存储方式，type唯一
	GetType() string
}

//单节点生效，多节点不采用此方式
type LocalStoreService struct {
	//此目录应该只能有读写权限
	Path         string
	StorePersist entity.StorePersist
	Type         string
}

func (localStoreService *LocalStoreService) createDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		logs.Error(err)
		return nil
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 600)
	}
	return err
}

func (localStoreService *LocalStoreService) Save(nameSpace string, content *[]byte, suffix string) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		logs.Error(err)
		return "", err
	}
	url := uid.String()
	err = localStoreService.createDirIfNotExist(localStoreService.Path + "/" + nameSpace)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	f, err := os.Create(localStoreService.Path + "/" + nameSpace + "/" + url)
	if f == nil || err != nil {
		logs.Error(err)
		return "", err
	}
	defer f.Close()
	_, err = f.Write(*content)
	now := time.Now()
	if err == nil {
		logs.Error(err)
		return url, localStoreService.StorePersist.Save(&entity.Store{Url: &url, StoreType: &localStoreService.Type, Suffix: suffix, CreateTime: &now})
	}
	return url, err
}
func (localStoreService *LocalStoreService) Get(nameSpace string, url string) ([]byte, string, error) {
	s, err := localStoreService.StorePersist.Get(url)
	if err != nil {
		logs.Error(err)
		return nil, "", err
	}
	if s == nil {
		return nil, "", nil
	}

	f, err := os.Open(localStoreService.Path + "/" + nameSpace + "/" + url)
	if f == nil || err != nil {
		logs.Error(err)
		return nil, "", err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logs.Error(err)
		return nil, "", err
	}
	return b, url + s.Suffix, nil
}

func (localStoreService *LocalStoreService) GetType() string {
	return localStoreService.Type
}
