package funcs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type Cache[T any] interface {
	Get(key string, data *T) (exists bool, err error)
	Set(key string, data T, duration time.Duration) error
}

func Remember[T any](key string, duration time.Duration, dst *T, fetchFunc func() (T, error), cache Cache[T]) error {
	md5Key := Md5Lower(key)
	exists, err := cache.Get(md5Key, dst)
	if err != nil {
		return err
	}
	if exists { // 正常取到直接返回
		return nil
	}
	*dst, err = fetchFunc()
	if err != nil {
		return err
	}
	err = cache.Set(key, *dst, duration)
	if err != nil {
		return err
	}
	return nil
}

func NewRedisCache[T any](client func() *redis.Client) Cache[T] {
	return RedisCache[T]{
		client: OnceValue(func() *redis.Client {
			return client()

		}),
	}
}

type RedisCache[T any] struct {
	client func() *redis.Client
}

func (r RedisCache[T]) Get(key string, data *T) (exists bool, err error) {
	ctx := context.Background()
	b, err := r.client().Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) { // 是redis.Nil 错误，屏蔽错误，exists 返回false
			return false, nil
		}
		return false, err
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (r RedisCache[T]) Set(key string, data T, duration time.Duration) (err error) {
	ctx := context.Background()
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = r.client().Set(ctx, key, b, duration).Result()
	if err != nil {
		return err
	}
	return nil
}

type MemeryCache[T any] struct {
	client func() *cache.Cache
}

func NewMemeryCache[T any](defaultExpiration time.Duration) Cache[T] {
	return MemeryCache[T]{
		client: OnceValue(func() *cache.Cache {
			return cache.New(defaultExpiration, 10*time.Minute)
		}),
	}
}

func (m MemeryCache[T]) Get(key string, dst *T) (exists bool, err error) {
	result, found := m.client().Get(key)
	if !found {
		return false, nil
	}
	*dst = result.(T)

	return true, nil
}

func (m MemeryCache[T]) Set(key string, data T, duration time.Duration) error {
	m.client().Set(key, data, duration)
	return nil
}

// 1秒缓存的内存缓存实例 常用于单次请求,某个接口、sql 结果的缓存，当流程函数中，多个组合调用同一个资源函数，但是个单元之间只能接受资源标识时，非常有用，1个接口通常需要在1s内执行完成

func MemeryCacheOneSecond[T any]() Cache[T] {
	return NewMemeryCache[T](1 * time.Second)
}
