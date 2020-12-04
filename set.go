package redis

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"strconv"
)

type ISetValue interface {
	ProtoBuf(value proto.Message) (reply interface{}, err error)
	String(value string) (reply interface{}, err error)
	Int(value int) (reply interface{}, err error)
	Int64(value int64) (reply interface{}, err error)
	Float64(value float64) (reply interface{}, err error)
	Json(value interface{}) (reply interface{}, err error)
	Bytes(value []byte) (reply interface{}, err error)
	Value(value interface{}) (reply interface{}, err error)
}
type SetValue struct {
	Conn redis.Conn
	Key  string
	Ex   int
}

func (s *SetValue) ProtoBuf(value proto.Message) (reply interface{}, err error) {
	defer s.Conn.Close()
	data, err := proto.Marshal(value)
	if err != nil {
		return nil, err
	}
	return s.Value(data)
}
func (s *SetValue) String(value string) (reply interface{}, err error) {
	defer s.Conn.Close()
	return s.Value(value)
}
func (s *SetValue) Json(value interface{}) (reply interface{}, err error) {
	defer s.Conn.Close()
	b, _ := json.Marshal(&value)
	return s.Value(b)
}
func (s *SetValue) Bytes(value []byte) (reply interface{}, err error) {
	defer s.Conn.Close()
	return s.Value(value)
}
func (s *SetValue) Value(value interface{}) (reply interface{}, err error) {
	defer s.Conn.Close()
	if s.Ex < 0 {
		return s.Conn.Do("SET", s.Key, value)
	}
	return s.Conn.Do("SET", s.Key, value, "EX", strconv.Itoa(s.Ex))

}

func (s *SetValue) Int(value int) (reply interface{}, err error) {
	defer s.Conn.Close()
	return s.Value(value)
}

func (s *SetValue) Int64(value int64) (reply interface{}, err error) {
	defer s.Conn.Close()
	return s.Value(value)
}
func (s *SetValue) Float64(value float64) (reply interface{}, err error) {
	defer s.Conn.Close()
	return s.Value(value)
}
