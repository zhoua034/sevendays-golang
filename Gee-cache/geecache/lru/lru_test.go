package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	//添加key1
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatal("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatal("cache miss key2 failed")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "v3"

	caps := len(k1 + k2 + v1 + v2)
	lru := New(int64(caps), nil)
	lru.Add(k1, String(v1))
	t.Log("添加" + k1)
	t.Log(lru.Len())
	lru.Add(k2, String(v2))
	t.Log("添加" + k2)
	t.Log(lru.Len())
	lru.Add(k3, String(v3))
	t.Log("添加" + k3)
	t.Log(lru.Len())
	//添加第三个元素超过内存，会移除第一个

	//验证
	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatal("lru RemoveOldest K1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	removedKeys := make([]string, 0)
	// 回掉函数，把淘汰的key添加到【】
	callback := func(key string, value Value) {
		removedKeys = append(removedKeys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))
	t.Log(removedKeys)
	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(expect, removedKeys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
