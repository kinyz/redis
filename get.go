package redis

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

type IGetValue interface {
	Value() (reply interface{}, err error)
	ProtoBuf(value proto.Message) error
	String() (string, error)
	Int() (int, error)
	Int64() (int64, error)
	Float64() (float64, error)
	Json(value interface{}) error
	Bytes() ([]byte, error)
}
type GetValue struct {
	Conn redis.Conn
	Key  string
}

func (g *GetValue) Value() (reply interface{}, err error) {
	defer g.Conn.Close()
	return g.Conn.Do("Get", g.Key)
}

func (g *GetValue) ProtoBuf(value proto.Message) error {
	bytes, err := redis.Bytes(g.Value())
	if err != nil {
		return err
	}
	return proto.Unmarshal(bytes, value)
}
func (g *GetValue) String() (string, error) {
	v, err := redis.String(g.Value())
	if err != nil {
		return "", err
	}
	return v, nil
}
func (g *GetValue) Int() (int, error) {
	v, err := redis.Int(g.Value())
	if err != nil {
		return 0, err
	}
	return v, nil
}
func (g *GetValue) Int64() (int64, error) {
	v, err := redis.Int64(g.Value())
	if err != nil {
		return 0, err
	}
	return v, nil
}
func (g *GetValue) Float64() (float64, error) {
	v, err := redis.Float64(g.Value())
	if err != nil {
		return 0, err
	}
	return v, nil
}
func (g *GetValue) Json(value interface{}) error {
	bytes, err := redis.Bytes(g.Value())
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}
func (g *GetValue) Bytes() ([]byte, error) {
	v, err := redis.Bytes(g.Value())
	if err != nil {
		return []byte(""), err
	}
	return v, nil
}
