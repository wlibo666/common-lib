package rediscluster

import (
	"strconv"
)

const (
	REDIS_CMD_CHAN_NUM = 10000

	CMD_INVALID = iota
	CMD_SETSTR
	CMD_SETNXSTR
	CMD_HSETSTR
	CMD_INCRBYSTR
)

type RedisCmd struct {
	Type int
	Cmd  []string
}

var (
	defChan = make(chan *RedisCmd, REDIS_CMD_CHAN_NUM)
)

func StartRedisChan() {
	for i := 0; i < defPoolSize; i++ {
		go func() {
			for cmd := range defChan {
				var err error = nil
				switch cmd.Type {
				case CMD_SETSTR:
					if len(cmd.Cmd) == 2 {
						_, err = Set(cmd.Cmd[0], cmd.Cmd[1], 0)
					} else if len(cmd.Cmd) == 3 {
						expire, _ := strconv.ParseInt(cmd.Cmd[2], 10, 64)
						_, err = Set(cmd.Cmd[0], cmd.Cmd[1], expire)
					}
				case CMD_SETNXSTR:
					if len(cmd.Cmd) == 2 {
						_, err = SetNx(cmd.Cmd[0], cmd.Cmd[1], 0)
					} else if len(cmd.Cmd) == 3 {
						expire, _ := strconv.ParseInt(cmd.Cmd[2], 10, 64)
						_, err = SetNx(cmd.Cmd[0], cmd.Cmd[1], expire)
					}
				case CMD_HSETSTR:
					if len(cmd.Cmd) == 3 {
						_, err = Hset(cmd.Cmd[0], cmd.Cmd[1], cmd.Cmd[2])
					}
				case CMD_INCRBYSTR:
					if len(cmd.Cmd) == 2 {
						cnt, _ := strconv.ParseInt(cmd.Cmd[1], 10, 64)
						_, err = IncrBy(cmd.Cmd[0], cnt)
					}
				default:
				}
				if err != nil {
					defChan <- cmd
				}
			}
		}()
	}
}

func newRedisCmd(Type int) *RedisCmd {
	return &RedisCmd{
		Type: Type,
	}
}

func SendRedisCmd(Type int, values ...string) {
	cmd := newRedisCmd(Type)
	cmd.Cmd = append(cmd.Cmd, values...)
	defChan <- cmd
}
