package geecache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	localCounts := make(map[string]int, len(db))
	for k, _ := range db {
		localCounts[k] = 0
	}
	geeCache := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		//回掉函数
		log.Printf("查询数据库，key:%s\n", key)
		if v, ok := db[key]; ok {
			//数据库有，记录从db查询次数
			localCounts[v] += 1
			return []byte(v), nil
		}
		//数据库中也没有，返回错误
		return nil, fmt.Errorf("%s 不存在！\n", key)
	}))

	//写入缓存
	for k, v := range db {
		//初始化缓存
		if view, err := geeCache.Get(k); err != nil || view.String() != v {
			log.Printf("failed to get key%s\n", k)
		}
		// 再次查询缓存
		if _, err := geeCache.Get(k); err != nil || localCounts[k] > 1 {
			log.Printf("cache miss %s\n", k)
		}
	}
}
