package framework

import (
	"errors"
	"fmt"
	"sync"
)

type Container interface {
	// 绑定一个服务提供者 如果凭证已经存在则替换
	Bind(provider ServiceProvider) error

	// 凭证是否已经绑定
	IsBind(key string) bool

	// 根据关键字获取服务
	Make(key string) (interface{}, error)

	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}

	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

type AnanContainer struct {
	Container

	// 存储注册服务的提供者
	providers map[string]ServiceProvider

	// 存储具体是实例
	instances map[string]interface{}

	// bind 和 make 防止并发操作
	lock sync.RWMutex
}

func NewAnanContainer() *AnanContainer {
	return &AnanContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (container *AnanContainer) PrintProviders() []string {
	ret := []string{}
	for _, provider := range container.providers {
		name := provider.Name()

		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

func (container *AnanContainer) Bind(provider ServiceProvider) error {
	container.lock.RLock()
	defer container.lock.RUnlock()

	key := provider.Name()

	container.providers[key] = provider //绑定服务提供者

	if !provider.IsDefer() { //不延迟实例化
		if err := provider.Boot(container); err != nil { //做一些初始化操作
			return err
		}

		params := provider.Params(container)
		method := provider.Register(container) //获取 实例化的方法
		ins, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		container.instances[key] = ins
	}
	return nil
}

func (container *AnanContainer) IsBind(key string) bool {
	return container.findServiceProvider(key) != nil
}

// 根据key找服务提供者
func (container *AnanContainer) findServiceProvider(key string) ServiceProvider {
	container.lock.RLock()
	defer container.lock.RUnlock()

	if sp, ok := container.providers[key]; ok {
		return sp
	}

	return nil
}

// 实例化服务提供者
func (container *AnanContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(container); err != nil {
		return nil, err
	}

	if params == nil {
		params = sp.Params(container)
	}

	method := sp.Register(container)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, nil
}

func (container *AnanContainer) Make(key string) (interface{}, error) {
	return container.make(key, nil, false)
}

func (container *AnanContainer) MustMake(key string) interface{} {
	serv, err := container.make(key, nil, false)
	if err != nil {
		panic(err)
	}

	return serv
}

func (container *AnanContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return container.make(key, params, true)
}

// 真正实例化一个服务
func (container *AnanContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	container.lock.RLock()
	defer container.lock.RUnlock()

	// 从提供者中寻找
	sp := container.findServiceProvider(key)

	if sp == nil {
		return nil, errors.New("contract:" + key + " is not found")
	}

	if forceNew { // 是否需要根据参数重新实例化
		return container.newInstance(sp, params)
	}

	// 查看实例中是否已经存在 存在的话返回
	if ins, ok := container.instances[key]; ok {
		return ins, nil
	}

	inst, err := container.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	container.instances[key] = inst
	return inst, nil
}

// 返回所有服务的名字
func (container *AnanContainer) NameList() []string {
	var ret = []string{}
	for _, provider := range container.providers {
		name := provider.Name()
		ret = append(ret, name)
	}

	return ret
}
