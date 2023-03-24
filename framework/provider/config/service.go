package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chsir-zy/anan/framework"
	"github.com/chsir-zy/anan/framework/contract"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
)

type AnanConfig struct {
	c        framework.Container
	folder   string // 文件夹
	keyDelim string //路径分隔符，默认是.
	lock     sync.RWMutex

	envMaps  map[string]string      // 所有的环境变量
	confMaps map[string]interface{} //配置文件
	confRaws map[string][]byte      //配置文件原始信息
}

// 读取某个配置文件
func (conf *AnanConfig) loadConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}

		// 对文本做环境变量的替换
		bf = replace(bf, conf.envMaps)
		// 解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, c); err != nil {
			return err
		}

		conf.confMaps[name] = c
		conf.confRaws[name] = bf

		// 读取app.path中的信息  更新app对应的folder
		if name == "app" && conf.c.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				appservice := conf.c.MustMake(contract.AppKey).(contract.App)
				appservice.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}

	return nil
}

// 删除文件
func (conf *AnanConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		delete(conf.confRaws, name)
		delete(conf.confMaps, name)
	}
	return nil
}

func NewAnanConfig(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
	}

	ananConf := &AnanConfig{
		c:        container,
		folder:   envFolder,
		envMaps:  envMaps,
		confMaps: map[string]interface{}{},
		confRaws: map[string][]byte{},
		keyDelim: ".",
		lock:     sync.RWMutex{},
	}

	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, file := range files {
		fileName := file.Name()
		err := ananConf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
		}
		continue
	}

	// 监控文件夹 看文件夹下面的文件是否有改动
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watch.Add(envFolder)
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watch.Events:
				{
					//判断事件发生的类型
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件：", ev.Name)
						ananConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
						ananConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						ananConf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error: ", err)
					return
				}
			}
		}
	}()

	return ananConf, nil
}

// replace 表示使用环境变量maps替换context中的env(xxx)的环境变量
func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}

	return content
}

// 查找某个路径下的配置项
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}

		switch next := next.(type) {
		case map[string]interface{}:
			return searchMap(next, path[1:])
		case map[interface{}]interface{}:
			return searchMap(cast.ToStringMap(next), path[1:])
		default:
			return nil
		}
	}

	return nil
}

func (conf *AnanConfig) find(key string) interface{} {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

func (conf *AnanConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

func (conf *AnanConfig) Get(key string) interface{} {
	return conf.find(key)
}

func (conf *AnanConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

func (conf *AnanConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

func (conf *AnanConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}
func (conf *AnanConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}
func (conf *AnanConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}
func (conf *AnanConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}
func (conf *AnanConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

func (conf *AnanConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

func (conf *AnanConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}
func (conf *AnanConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

func (conf *AnanConfig) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(conf.find(key))
}
