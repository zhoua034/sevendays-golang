package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nowBytes  int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

// entry List中节点的类型
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is a Constructor of Cache
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		cache:     make(map[string]*list.Element),
		ll:        list.New(),
		maxBytes:  maxBytes,
		OnEvicted: onEvicted,
	}
}

// Get look for the key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 如果元素存在，移动到队头，并返回
	if elem, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elem)
		elem_kv := elem.Value.(*entry)
		return elem_kv.value, true
	}
	return
}

// RemoveOldest remove the last element
func (c *Cache) RemoveOldest() {
	//拿到队尾元素
	ele := c.ll.Back()
	if ele != nil {
		//删除队尾元素
		c.ll.Remove(ele)
		//删除cache中
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		//更新内存
		c.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//回掉函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add to add a element to cache
func (c *Cache) Add(key string, value Value) {
	//是否存在
	if elem, ok := c.cache[key]; ok {
		//移动队头
		c.ll.MoveToFront(elem)

		kv := elem.Value.(*entry)
		//更新内存
		c.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		//更新值
		kv.value = value
	} else {
		//添加新值
		newElem := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = newElem
		//更新内存
		c.nowBytes += int64(len(key)) + int64(value.Len())
	}
	//如果超过了max，就一直removeOldest
	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.RemoveOldest()
	}
}

// Len for test
func (c *Cache) Len() int {
	return c.ll.Len()
}
