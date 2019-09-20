package xid

import (
	"errors"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func lock(client *redis.Client, lockKey string, timeoutMs, retryTimes int) (UnlockFunc, error) {
	resp, err := lockRetry(client, lockKey, timeoutMs, retryTimes)
	if !resp {
		return nil, errors.New("lock error: " + err)
	}
	return func() {
		delResp := client.Del(lockKey)
		_, _ = delResp.Result()
	}, nil
}

func expiredTime(timeoutMs int) (int64, int64) {
	now := time.Now()
	return now.UnixNano(), now.Add(time.Duration(timeoutMs * 1000000)).UnixNano()
}

// 重试 retry_times 次
func lockRetry(rds redis.Cmdable, key string, timeoutMs, retryTimes int) (bool, string) {
	for i := 0; i < retryTimes; i++ {
		now, ex := expiredTime(timeoutMs)
		setNxCmd := rds.SetNX(key, ex, 0)
		if setNxCmd.Val() {
			return true, strconv.FormatInt(ex, 10)
		}
		getCmd := rds.Get(key)
		if getCmd.Val() == "" {
			getSetCmd := rds.GetSet(key, ex)
			if getSetCmd.Val() == "" {
				return true, strconv.FormatInt(ex, 10)
			}
		} else {
			prevTime, err := getCmd.Int64()
			if err != nil {
				log.Println("get key int64 err:", err)
			} else {
				// 已经过期，可以尝试获得锁了
				if now > prevTime {
					getSetCmd2 := rds.GetSet(key, ex)
					if getSetCmd2.Val() == getCmd.Val() {
						return true, strconv.FormatInt(ex, 10)
					}
				}
			}
		}
		wait := rand.Int63n(100) + 300
		time.Sleep(time.Millisecond * time.Duration(wait))
	}
	return false, "0"
}
