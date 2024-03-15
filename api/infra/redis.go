package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	conn *redis.Client
}

var redisDbLock *Redis
var once sync.Once

const DB_LOCK_KEY_NAME = "lock"

func GetInstance(redisUrl string) *Redis {
	if redisDbLock == nil {
		once.Do(
			func() {

				r := redis.NewClient(&redis.Options{
					Addr:     redisUrl,
					Password: "",
					DB:       0,
				})

				redis := &Redis{
					conn: r,
				}

				ping, err := r.Ping(context.Background()).Result()

				if err != nil {
					fmt.Printf("Error connecting to Redis: %s \n", err)
				} else {
					// set default value for lock
					r.Conn().Set(context.Background(), DB_LOCK_KEY_NAME, "0", 0)
					fmt.Printf("Connected to Redis successfully. Ping: %s \n", ping)
				}

				redisDbLock = redis
			},
		)
	}

	return redisDbLock
}

func (r *Redis) GetDbLock(ctx context.Context) (string, error) {
	val, err := r.conn.Get(ctx, DB_LOCK_KEY_NAME).Result()

	if err != nil {
		fmt.Printf("Error getting lock: %s \n", err)
		return "", err
	}

	return val, nil
}

func (r *Redis) LockDb(ctx context.Context) error {
	err := r.conn.Set(ctx, DB_LOCK_KEY_NAME, "1", 0).Err()
	if err != nil {
		return err
	}

	return nil

}

func (r *Redis) UnlockDb(ctx context.Context) error {
	err := r.conn.Set(ctx, DB_LOCK_KEY_NAME, "0", 0).Err()
	if err != nil {
		return err
	}

	return nil

}
