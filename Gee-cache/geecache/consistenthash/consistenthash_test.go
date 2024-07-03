package consistenthash

import (
	"log"
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	//添加真实节点 2,4,6 会添加到hashRing中：2，4，6，12，14，16，22，24，26
	hash.Add("2", "4", "6")

	//测试请求 2,11,27 ,应命中虚拟节点 2,12,,24,2 》真实节点为 2，2，4,2
	testcase := []string{"2", "11", "23", "27"}
	for _, key := range testcase {
		vrNode := hash.Get(key)
		log.Printf("%s命中的真实节点为%s\n", key, vrNode)
	}
	log.Println("》》》》》《《《《")
	//添加真实节点8
	//hashRing 》 2，4，6，8，12，14，16，18，22，24，26，28
	hash.Add("8")
	//再次测试6，会命中虚拟节点28》真实节点8
	testKey := "27"
	get := hash.Get(testKey)
	log.Printf("再次测试%s命中的真实节点为%s\n", testKey, get)
}
