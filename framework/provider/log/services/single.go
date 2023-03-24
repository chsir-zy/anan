package services

import (
	"os"
	"path/filepath"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/chsir-zy/anan/framework/util"

	"github.com/pkg/errors"
)

type AnanSingleLog struct {
	AnanLog

	folder string
	file   string
}

func NewAnanSingleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &AnanSingleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	appService := c.MustMake(contract.AppKey).(contract.App)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	folder := appService.LogFolder()
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}

	log.folder = folder
	if !util.Exists(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	log.file = "anan.log"
	if configService.IsExist("log.file") {
		log.file = configService.GetString("log.file")
	}

	fd, err := os.OpenFile(filepath.Join(folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file err")
	}

	// 将内容输出到控制台
	log.SetOutput(fd)
	log.c = c
	return log, nil
}
