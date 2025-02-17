package dbsaver

import (
	"io"

	"github.com/admpub/log"
	modelFile "github.com/admpub/nging/v3/application/model/file"
	"github.com/fatih/color"
	uploadClient "github.com/webx-top/client/upload"
)

type (
	DBSaver func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error
)

var (
	dbSavers = map[string]DBSaver{}
	Default  = func(fileM *modelFile.File, result *uploadClient.Result, reader io.Reader) error {
		return fileM.Add(reader)
	}
)

// Register 注册数据保存处理函数
// key table.field
func Register(key string, dbsaver DBSaver) {
	dbSavers[key] = dbsaver
	log.Info(color.YellowString(`dbsaver.register:`), key)
}

func Unregister(keys ...string) {
	for _, key := range keys {
		_, ok := dbSavers[key]
		if ok {
			delete(dbSavers, key)
		}
	}
}

func Get(key string, defaults ...string) DBSaver {
	if dbsaver, ok := dbSavers[key]; ok {
		return dbsaver
	}
	if len(defaults) > 0 {
		return Get(defaults[0])
	}
	return Default
}
