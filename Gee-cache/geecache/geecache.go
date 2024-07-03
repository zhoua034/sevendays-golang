package geecache

import (
	"fmt"
	pb "geecache/geecachepb"
	"geecache/singleflight"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	// 用于寻找远程节点的PeerPicker 》 HTTPPool实现：
	// 既能提供http服务，又能找到其他节点
	peer PeerPicker
	// single flight 防止缓存击穿，合并相同大量的请求
	loader *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create new instance of group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func (g *Group) RegisterPeers(peer PeerPicker) {
	if g.peer != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peer = peer
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	group := groups[name]
	mu.RUnlock()
	return group
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is nil")
	}
	//从缓存中查找
	if view, ok := g.mainCache.get(key); ok {
		log.Printf("本机 [Cache] hit %s!\n", key)
		return view, nil
	}

	//缓存未命中从回掉函数查询
	return g.load(key)

}

func (g *Group) load(key string) (value ByteView, err error) {

	//使用single flight Do包裹起来，确保大量相同的key只请求一次
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peer != nil {
			if peer, ok := g.peer.PickPeer(key); ok {
				// 若peer为远程节点，则从远程peer获取
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		//若一致性hash之后不是远程节点，则从回掉函数取（从数据库取）
		return g.loadLocally(key)
	})
	if err == nil {
		return view.(ByteView), nil
	}

	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	// 从远程节点获取
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

func (g *Group) loadLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{
		b: cloneBytes(bytes),
	}
	// 更新缓存
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
