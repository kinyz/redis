package redis

import "github.com/garyburd/redigo/redis"

type IValue interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
	Check(key string) bool
	Del(key string) (reply interface{}, err error)
	Get(key string) IGetValue
	Set(key string, timeOut ...int) ISetValue
}
type Value struct {
	Conn redis.Conn
}

func (v *Value) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	defer v.Conn.Close()
	return v.Conn.Do(commandName, args)
}

func (v *Value) Check(key string) bool {

	defer v.Conn.Close()
	_, err := v.Conn.Do("EXISTS", key)
	if err != nil {
		return false
	}
	return true
}
func (v *Value) Del(key string) (reply interface{}, err error) {
	defer v.Conn.Close()
	return v.Conn.Do("DEL", key, "SEX")
}

func (v *Value) Get(key string) IGetValue {
	return &GetValue{Conn: v.Conn, Key: key}
}

func (v *Value) Set(key string, timeOut ...int) ISetValue {
	if len(timeOut) < 1 {
		return &SetValue{Conn: v.Conn, Key: key, Ex: -1}
	}
	return &SetValue{Conn: v.Conn, Key: key, Ex: timeOut[0]}
}
