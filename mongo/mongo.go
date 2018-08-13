package mongo

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

// mongo连接配置
type MongoInfo struct {
	Name     string
	Addr     string
	User     string
	Pwd      string
	AuthDb   string
	PoolSize int
	Timeout  int
	Session  *mgo.Session

	del bool
}

var (
	// 可能需要连接多个MongoDB集群
	GMongo      map[string]*MongoInfo = make(map[string]*MongoInfo)
	GMongoMutex *sync.Mutex           = &sync.Mutex{}
)

// 添加MongoDB连接配置
func AddConf(name, addr, user, pwd, authDb string, poolSize, timeout int) error {
	GMongoMutex.Lock()
	defer GMongoMutex.Unlock()
	_, ok := GMongo[name]
	if !ok {
		info := &MongoInfo{
			Addr:     addr,
			User:     user,
			Pwd:      pwd,
			AuthDb:   authDb,
			PoolSize: poolSize,
			Timeout:  timeout,
		}
		GMongo[name] = info
		info.check()
	} else {
		return fmt.Errorf("mongo [%s] exists", name)
	}
	return nil
}

func DelConf(name string) {
	GMongoMutex.Lock()
	defer GMongoMutex.Unlock()
	info, ok := GMongo[name]
	if ok {
		info.del = true
		if info.Session != nil {
			info.Session.Close()
		}
		time.Sleep(2 * time.Second)
		delete(GMongo, name)
	}
}

func Close(name string) {
	GMongoMutex.Lock()
	defer GMongoMutex.Unlock()
	mInfo, ok := GMongo[name]
	if ok {
		mInfo.Session.Close()
		mInfo.Session = nil
	}
}

// 连接到MongoDB
func (mInfo *MongoInfo) connect() error {
	var err error
	mInfo.Session, err = mgo.DialWithTimeout(mInfo.Addr, time.Duration(mInfo.Timeout)*time.Second)
	if err != nil {
		return err
	}
	if mInfo.User != "" && mInfo.Pwd != "" && mInfo.AuthDb != "" {
		auth := &mgo.Credential{
			Username: mInfo.User,
			Password: mInfo.Pwd,
			Source:   mInfo.AuthDb,
		}
		err = mInfo.Session.Login(auth)
		if err != nil {
			mInfo.Session.Close()
			return err
		}
	}
	if mInfo.PoolSize > 0 {
		mInfo.Session.SetPoolLimit(mInfo.PoolSize)
	}
	if mInfo.Timeout > 0 {
		mInfo.Session.SetSocketTimeout(time.Duration(mInfo.Timeout) * time.Second)
		mInfo.Session.SetSyncTimeout(time.Duration(mInfo.Timeout) * time.Second)
	}

	mInfo.Session.SetMode(mgo.Nearest, true)
	return nil
}

// 检测到MongoDB的连接是否中断，若已断开则重连
func (mInfo *MongoInfo) check() {
	go func(mInfo *MongoInfo) {
		for {
			if mInfo.del {
				return
			}
			if mInfo.Session != nil {
				if mInfo.Session.Ping() != nil {
					mInfo.connect()
				}
			}
			time.Sleep(time.Second)
		}
	}(mInfo)
}

// 根据mongo配置名称获取一个连接
func GetSession(name string) (*mgo.Session, error) {
	mInfo, ok := GMongo[name]
	if !ok {
		return nil, fmt.Errorf("not found mongo by name:%s", name)
	}
	if mInfo.Session == nil {
		err := mInfo.connect()
		if err != nil {
			return nil, err
		}
	}
	return mInfo.Session.Clone(), nil
}
