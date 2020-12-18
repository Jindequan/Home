package redis

import (
	"bigdata_permission/conf"
	"github.com/go-redis/redis"
	"time"
)

var (
	client    *redis.Client
	keyPrefix = "web_leads_backend::"
)

func getClient() (bool, *redis.Client) {
	options := &redis.Options{
		Addr:     conf.BaseConf.RedisConfig.Addr,
		Password: conf.BaseConf.RedisConfig.Password,
		DB:       conf.BaseConf.RedisConfig.Db,
	}
	client = redis.NewClient(options)
	res := client.Ping()
	if res == nil {
		return false, &redis.Client{}
	}
	_, err := res.Result()
	if err != nil {
		return false, &redis.Client{}
	}
	return true, client
}

func checkConnection() bool {
	if client == nil {
		return false
	}
	res := client.Ping()
	if res == nil {
		return false
	}
	_, err := res.Result()
	if err != nil {
		return false
	}
	return true
}

func initClient() bool {
	if checkConnection() {
		return true
	}
	success, _ := getClient()
	return success
}

func genKey(k string) string {
	return keyPrefix + k
}

func genTimeout(t int) time.Duration {
	if t == 0 {
		//不允许永久key的存在
		t = 3600
	}
	return time.Duration(int64(t)) * time.Second
}

func Set(k, v string, timeout int) bool {
	if !initClient() {
		return false
	}

	err := client.Set(genKey(k), v, genTimeout(timeout)).Err()
	return err == nil
}

func Get(k string) (bool, string) {
	if !initClient() {
		return false, ""
	}

	res, err := client.Get(genKey(k)).Result()
	return err == nil, res
}

func HSet(k, f, v string) bool {
	if !initClient() {
		return false
	}

	err := client.HSet(genKey(k), f, v).Err()
	return err == nil
}

func HGet(k, f string) (bool, string) {
	if !initClient() {
		return false, ""
	}

	res, err := client.HGet(genKey(k), f).Result()
	return err == nil, res
}

func Incr(k string) (bool, int) {
	if !initClient() {
		return false, 0
	}

	res, err := client.Incr(genKey(k)).Result()
	return err == nil, int(res)
}

func Decr(k string) (bool, int) {
	if !initClient() {
		return false, 0
	}

	res, err := client.Decr(genKey(k)).Result()
	return err == nil, int(res)
}

func SetNX(k string, timeout int) bool {
	if !initClient() {
		return false
	}

	lockSuccess, err := client.SetNX(genKey(k), 1, genTimeout(timeout)).Result()
	return err == nil && lockSuccess
}

func Del(k string) bool {
	if !initClient() {
		return false
	}

	delNum, err := client.Del(genKey(k)).Result()
	return err == nil && delNum > 0
}

func Ttl(k string) (bool, int) {
	if !initClient() {
		return false, 0
	}

	ttl, err := client.TTL(genKey(k)).Result()
	if err == nil {
		return true, int(ttl.Microseconds()) / 1000
	}
	return false, 0
}