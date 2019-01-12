package redis

import (
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
	once sync.Once
)

func ConnectInit(addr, passwd, db string) error {
	once.Do(func() {
		pool = &redis.Pool{
			MaxIdle:     2,
			MaxActive:   20,
			IdleTimeout: 180 * time.Second,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", addr)
				if err != nil {
					return nil, err
				}
				if strings.TrimSpace(passwd) != "" {
					if _, err := conn.Do("AUTH", strings.TrimSpace(passwd)); err != nil {
						conn.Close()
						return nil, err
					}
				}
				conn.Do("SELECT", db)
				return conn, nil
			},
		}
	})
	return nil
}

func GetValue(key string) (interface{}, error) {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	return redis_cli.Do("GET", key)
}

func SetValue(key string, value interface{}) error {
	return SetValueAndExpire(key, value, 300)
}

func SetValueAndExpire(key string, value interface{}, expire int) error {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	_, err := redis_cli.Do("SET", key, value)
	if err != nil {
		return err
	}
	if expire == 0 {
		_, err = redis_cli.Do("EXPIRE", key, 300)
	} else {
		_, err = redis_cli.Do("EXPIRE", key, expire)
	}
	return err
}

func SetValueNoExpire(key string, value interface{}) error {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	_, err := redis_cli.Do("SET", key, value)
	return err
}

//设置key，仅当key不存在时成功，同时设置过期时间，def为"px"表示毫秒，"ex"表示秒
//设置成功返回“OK”，否则返回nil
func SetNxKeyAndExpire(key string, value interface{}, def string, expire int) (interface{}, error) {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	return redis_cli.Do("SET", key, value, def, expire, "nx")
}

func SetExpire(key string, expire ...int) error {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	tExpire := 300
	if len(expire) > 0 {
		tExpire = expire[0]
	}
	_, err := redis_cli.Do("EXPIRE", key, tExpire)
	return err
}

func Delete(key string) error {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	_, err := redis_cli.Do("DEL", key)
	return err
}

//设置redis中的set中的值
func SetSetValue(key string, value interface{}) error {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	_, err := redis_cli.Do("SADD", key, value)
	if err != nil {
		return err
	}
	_, err = redis_cli.Do("EXPIRE", key, 300)
	return err
}

//判断某个值是否在set中
func IsSetMember(key string, value interface{}) bool {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	exist, err := redis.Bool(redis_cli.Do("SISMEMBER", key, value))
	if err != nil {
		return false
	}
	return exist
}

//获取set集合中的所有值
func GetSetAllValue(key string) (interface{}, error) {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	return redis_cli.Do("SMEMBERS", key)
}

func ExistKey(key string) bool {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	exist, err := redis.Bool(redis_cli.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exist
}

//获取集合中元素的个数
func GetLenOfSet(key string) (int, error) {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	return redis.Int(redis_cli.Do("scard", key))
}

//获取链表中元素的个数
func GetLenOfList(key string) (int, error) {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	return redis.Int(redis_cli.Do("llen", key))
}

//从链表尾处插入元素，用户依据相应场景设置自身的过期时间
func PushElementWithTail(key string, value interface{}, expire int) error {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	_, err := redis_cli.Do("rpush", key, value)
	if err != nil {
		return err
	}
	if expire == 0 {
		_, err = redis_cli.Do("EXPIRE", key, 300)
	} else {
		/*expire_time,_:=strconv.Atoi(expire)*/
		_, err = redis_cli.Do("EXPIRE", key, expire)
	}
	return err
}

//从链表头位置删除元素
func PopElementFromHead(key string) (interface{}, error) {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	return redis_cli.Do("lpop", key)
}

func GetAllElementsFromList(key string) (interface{}, error) {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	return redis_cli.Do("lrange", key, 0, -1)
}

//按照用户所设定的过期时间进行处理
func SetValueBasedOnExpire(key string, value interface{}, expire int) error {
	return SetValueAndExpire(key, value, expire)
}

//按照用户所设定的过期时间，对集合中的元素进行处理
func SetSetValueBasedOnExpire(key string, value interface{}, expire int) error {
	redis_cli := pool.Get()
	defer redis_cli.Close()
	_, err := redis_cli.Do("SADD", key, value)
	if err != nil {
		return err
	}
	_, err = redis_cli.Do("EXPIRE", key, expire)
	return err
}
