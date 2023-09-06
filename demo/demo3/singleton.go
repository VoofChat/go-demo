package demo3

import "sync"

// 方法一: 懒汉式
// 懒汉式是一种常见的单例模式实现方式，其特点是在首次使用时创建单例实例。
/*
var (
	instance *Singleton
	once     sync.Once
)

type Singleton struct {
}

func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})

	return instance
}
*/

//方法二：饿汉式
//饿汉式是另一种常见的单例模式实现方式，其特点是在系统启动时即创建单例实例，当调用时直接返回该实例。
/*

var instance *Singleton = &Singleton{}

type Singleton struct {
}

func GetInstance() *Singleton {
	return instance
}
*/

// 方法三：双重检查锁定
// 双重检查锁定是一种在多线程环境下使用的单例模式实现方式，其特点是先检查是否已经有实例，如果没有则进入同步代码块进行创建。
var (
	instance *Singleton
	mu       sync.Mutex
)

type Singleton struct {
}

func GetInstance() *Singleton {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = &Singleton{}
		}
	}
	return instance
}
