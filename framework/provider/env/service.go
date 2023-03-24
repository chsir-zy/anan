package env

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/chsir-zy/anan/framework/contract"
)

// Env的具体实现
type AnanEnv struct {
	floder string            // 代表.env所在的具体目录
	maps   map[string]string // 保持所有的环境变量
}

// NewAnanEnv 有一个参数，.env文件所在的目录
// example: NewAnanEnv("/envfolder/") 会读取文件: /envfolder/.env
// .env的文件格式 FOO_ENV=BAR
func NewAnanEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("params error")
	}

	// 读取配置文件目录
	floder := params[0].(string)

	// 实例化
	ananEnv := &AnanEnv{
		floder: floder,
		// 默认为开发环境
		maps: map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	file := path.Join(floder, ".env")
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		// 读取文件
		br := bufio.NewReader(fi)
		for {
			// 按照行读取
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}

			s := bytes.SplitN(line, []byte{'='}, 2)
			// 如果不符合规范 过滤掉
			if len(s) < 2 {
				continue
			}

			key := string(s[0])
			val := string(s[1])
			ananEnv.maps[key] = val
		}
	}

	// 获取当前的环境变量，并且覆盖.env文件的配置
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}

		ananEnv.maps[pair[0]] = pair[1]
	}

	return ananEnv, nil
}

func (a *AnanEnv) AppEnv() string {
	return a.Get("APP_ENV")
}

func (a *AnanEnv) IsExist(key string) bool {
	_, ok := a.maps[key]
	return ok
}

func (a *AnanEnv) All() map[string]string {
	return a.maps
}

func (a *AnanEnv) Get(key string) string {
	if val, ok := a.maps[key]; ok {
		return val
	}

	return ""
}
