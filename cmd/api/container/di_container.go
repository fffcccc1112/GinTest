package container

import (
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"test/pkg/logger"
)

// 实现一个DI
type Destoryable interface {
	Destroy() error
}
type FactoryFunc func() (interface{}, error)

// 单例模式
// 工厂模式
// 销毁池
type DIContainer struct {
	instances       map[reflect.Type]interface{}
	instanceFactory map[reflect.Type]FactoryFunc
	mu              sync.Mutex
}

// 初始化
func NewDIContainer() *DIContainer {
	return &DIContainer{
		instances:       make(map[reflect.Type]interface{}),
		instanceFactory: make(map[reflect.Type]FactoryFunc),
		mu:              sync.Mutex{},
	}
}

// inject container
func (c *DIContainer) RegisterInstance(t reflect.Type) (interface{}, error) {
	//从工厂模式获取方法
	//先查看单例模式有无实例

	c.mu.Lock()
	defer c.mu.Unlock()
	instance, exists := c.instances[t]
	if exists {
		return instance, nil
	}
	factoryFunc, exists := c.instanceFactory[t]
	if exists {
		var err error
		c.instances[t], err = factoryFunc()
		if err != nil {
			return nil, fmt.Errorf("单例模式实例化失败!%s", err)
		}
	} else {
		logger.Error("单例模式未找到注册的工厂方法")
		return nil, fmt.Errorf("单例模式未找到注册的工厂方法")
	}
	return c.instances[t], nil
}
func (c *DIContainer) Register(t reflect.Type, f FactoryFunc) bool {
	//查看是否已经注册
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exist := c.instanceFactory[t]
	if exist {
		logger.Info("工厂方法已经注册,无需注册")
		return true
	}
	//注册
	c.instanceFactory[t] = f
	return true
}

// TODO：销毁时
func (c *DIContainer) DestoryInstance(t reflect.Type) error {
	c.mu.Lock()
	defer c.mu.Unlock() //销毁函数销毁
	if destoryable, ok := c.instances[t].(Destoryable); ok {
		//销毁
		err := destoryable.Destroy()
		if err != nil {
			logger.Error("销毁失败!", zap.String("类型", t.String()))
			return err
		}
	}

	delete(c.instances, t)
	return nil
}

// 封装自动注册工厂模式,
func RegisterFactoryAndInstance[T any, K any](c *DIContainer, t reflect.Type, preInstance K,
	f func(pre K) (*T, error)) (*T, error) {
	//注册
	c.Register(t, func() (interface{}, error) {
		instance, err := f(preInstance)
		if err != nil {
			return nil, fmt.Errorf("容器注册失败!")
		}
		return instance, nil
	},
	)
	//注册之后单例模式获取实例
	instance, err := c.RegisterInstance(t)
	if err != nil {
		return nil, err
	}

	result, ok := instance.(*T)
	if !ok {
		return nil, fmt.Errorf("类型断言失败：实例实际类型是 %T，无法转为目标类型 %T", instance, *new(T))
	}

	return result, nil
}
func RegisterFactoryAndBased[T any](c *DIContainer, t reflect.Type,
	f func() (*T, error)) (*T, error) {
	//注册

	c.Register(t, func() (interface{}, error) {
		instance, err := f()
		if err != nil {
			return nil, fmt.Errorf("容器注册失败!")
		}
		return instance, nil
	},
	)
	//注册之后单例模式获取实例
	instance, err := c.RegisterInstance(t)
	if err != nil {
		return nil, err
	}

	result, ok := instance.(*T)
	if !ok {
		return nil, fmt.Errorf("类型断言失败：实例实际类型是 %T，无法转为目标类型 %T", instance, *new(T))
	}

	return result, nil
}
