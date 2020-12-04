package redis

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type IPool interface {
	/*
		获取redis master pool
	*/
	Master() *redis.Pool
	/*
		获取redis slave pool
	*/
	Slave() *redis.Pool
	/*
		获取一个redis master conn
	*/
	GetMasterConn() IValue
	/*
		获取一个redis slave conn
	*/
	GetSlaveConn() IValue
	/*
		获取一个原始的redis master conn
	*/
	GetRawMasterConn() redis.Conn
	/*
		获取一个原始的redis slave conn
	*/
	GetRawSlaveConn() redis.Conn
	/*
		关闭redis master pool
	*/
	CloseMaster() error
	/*
		关闭redis slave pool
	*/
	CloseSlave() error
	/*
		分布式锁
	*/
	Lock() ILock
}
type Pool struct {
	MasterPool *redis.Pool
	SlavePool  *redis.Pool
	conf       *RedisConfig
}

func NewRedisPool(file string) IPool {
	pool := &Pool{}
	pool.conf = NewRedisConfig(file)

	pool.MasterPool = &redis.Pool{
		MaxIdle:     pool.conf.Master.MaxIdle,
		MaxActive:   pool.conf.Master.MaxActive,
		IdleTimeout: 100,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", pool.conf.GetMasterAddr(), redis.DialPassword(pool.conf.Master.PassWord))
			if err != nil {
				log.Println("Redis 连接失败")
				//panic(err)
				return nil, err
			}
			return conn, nil
		},
	}

	if pool.conf.Slave.Enable {
		pool.SlavePool = &redis.Pool{
			MaxIdle:     pool.conf.Slave.MaxIdle,
			MaxActive:   pool.conf.Slave.MaxActive,
			IdleTimeout: 100,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", pool.conf.GetSlaveAddr(), redis.DialPassword(pool.conf.Slave.PassWord))
				if err != nil {
					log.Println("Redis 连接失败", err)
					//panic(err)
					return nil, err
				}
				return conn, nil
			},
		}
	}

	return pool
}

func (p *Pool) Master() *redis.Pool {
	return p.MasterPool
}

func (p *Pool) Slave() *redis.Pool {
	return p.SlavePool
}

func (p *Pool) Lock() ILock {
	return &Lock{MasterConn: p.GetMasterConn(), SlaveConn: p.GetSlaveConn()}
}

// 获取redis Master conn
func (p *Pool) GetMasterConn() IValue {
	return &Value{Conn: p.Master().Get()}
}

// 获取redis Slave conn
func (p *Pool) GetSlaveConn() IValue {
	if !p.conf.Slave.Enable {
		return &Value{Conn: p.Master().Get()}
	}
	return &Value{Conn: p.Slave().Get()}
}

// 获取redis Master conn
func (p *Pool) GetRawMasterConn() redis.Conn {
	return p.Master().Get()
}

// 获取redis Slave conn
func (p *Pool) GetRawSlaveConn() redis.Conn {
	if !p.conf.Slave.Enable {
		return p.Master().Get()
	}
	return p.Slave().Get()
}

func (p *Pool) CloseMaster() error {
	return p.Master().Close()
}
func (p *Pool) CloseSlave() error {
	return p.Slave().Close()
}
