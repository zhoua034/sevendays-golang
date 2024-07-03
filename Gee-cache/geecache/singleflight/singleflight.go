package singleflight

import (
	"sync"
)

// call 正在进行中的请求
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 主结构管理所有key的请求call
type Group struct {
	mu sync.Mutex //protect the m
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock() //上锁
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		//已经有call正在执行
		//log.Printf("%s等待结果", key)
		g.mu.Unlock()
		//进入等待结果
		c.wg.Wait()
		return c.val, c.err
	}
	//否则绑定m
	c := new(call)
	c.wg.Add(1)
	//注册到g.m中
	g.m[key] = c
	g.mu.Unlock() //解锁

	//执行Do过程，拿到结果
	c.val, c.err = fn()
	c.wg.Done() //通知已进行完毕

	//删除g.m中的key
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}
