package redis

import (
	"errors"
	"log"
)

type ILock interface {
	/*
		锁
		参数列表:
		key string
		time int 秒
		value 所有类型
		返回值:
		error为空时获取锁成功 否则失败
	*/
	Lock(key string, time int, value ...interface{}) error
	/*
		解锁
		参数列表:
		key string
		返回值:
		error为空时解锁成功 否则失败
	*/
	UnLock(Key string) error

	/*
		查询锁
		参数列表:
		key string
		返回值:
		true 存在
		false 失败
	*/
	Check(Key string) bool
}

var lockName = "lock-" //设置redis里的前缀

type Lock struct {
	MasterConn IValue
	SlaveConn  IValue
	Namespace  string
}

func (l *Lock) Lock(key string, time int, value ...interface{}) error {
	if l.Check(key) {
		return errors.New("key already exists")
	}
	log.Println(value)
	_, err := l.MasterConn.Set(lockName+key, time).Value(value)
	if err != nil {
		return err
	}
	return nil

}
func (l *Lock) UnLock(Key string) error {

	_, err := l.MasterConn.Del(lockName + Key)
	if err != nil {
		log.Println("UnLock ", Key, "Error :", err)
		return err
	}
	return nil
}

func (l *Lock) Check(Key string) bool {
	v, err := l.SlaveConn.Get(lockName + Key).Value()
	if err != nil {
		return false
	}
	if v != nil {
		return true
	}
	return false
}
