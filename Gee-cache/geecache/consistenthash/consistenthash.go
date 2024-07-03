package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// Map is the consistent hash struct
type Map struct {
	hash      Hash           // hash算法
	vrNums    int            //虚拟节点个数
	hashRing  []int          //hash值环 Sorted
	vrHashMap map[int]string //虚拟节点和真实节点的映射
}

// New creates a Map instance
func New(vrNums int, fn Hash) *Map {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &Map{
		hash:      fn,
		vrNums:    vrNums,
		vrHashMap: make(map[int]string),
	}
}

// Add 添加一个或多个真实节点
func (m *Map) Add(keys ...string) {
	// 每个真实节点创建vrNums个虚拟节点 加入到hash环中
	for _, key := range keys {
		for i := 0; i < m.vrNums; i++ {
			//对节点拼接虚拟节点,并计算出has值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//把虚拟节点hash添加到hashRing中
			m.hashRing = append(m.hashRing, hash)
			//增加映射关系
			m.vrHashMap[hash] = key
		}
	}
	//环进行排序
	sort.Ints(m.hashRing)
}

func (m *Map) Get(key string) string {
	//hashRing为空
	if len(m.hashRing) == 0 {
		return ""
	}
	//计算key的hash值
	hash := int(m.hash([]byte(key)))
	//从hashRing中查找第一个大于key's hash的虚拟节点
	idx := sort.Search(len(m.hashRing), func(i int) bool {
		return m.hashRing[i] >= hash
	})
	//用%，因为如果是没有比key's hash大的 idx为len(m.hashRing)
	//也就是该使用第0个位置的虚拟节点
	k := m.hashRing[idx%len(m.hashRing)]
	return m.vrHashMap[k]
}
