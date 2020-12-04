package redis

import (
	"log"
	"testing"
)

func TestRun(t *testing.T) {

	pool := NewRedisPool("conf/redis.yaml")
	defer pool.CloseMaster()
	defer pool.CloseSlave()
	pool.GetMasterConn().Set("testKey", 100).String("")
	info, err := pool.GetMasterConn().Get("testKey").Value()

	pool.Lock().Lock("uuid", 20, "key1", 1)
	pool.Lock().UnLock("uuid")
	if err := pool.Lock().Lock("uuid", 20); err != nil {
		log.Println(err)
	}
	log.Print(info, err)
	select {}
}
