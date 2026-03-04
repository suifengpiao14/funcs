package funcs

import (
	"io"
	"time"
)

type BackoffDuration []time.Duration

var DefaultBackoffDuration = BackoffDuration{1 * time.Second, 2 * time.Second, 4 * time.Second}

// ExponentialBackoff 指数退避重试函数
func ExponentialBackoff(backoffDuration BackoffDuration, needReDo func(err error) bool, dofn func() error) (err error) {
	if len(backoffDuration) == 0 { // 1. 没有退避时间，直接执行一次
		return dofn()
	}
	for _, backoff := range backoffDuration {
		err = dofn()
		if err == nil { // 2. 执行成功，直接返回
			return err
		}
		if !needReDo(err) { // 3. 有错误，但是不需要重试，直接返回
			return err
		}
		// 4. 需要重试，等待一段时间再执行
		if backoff > 0 {
			time.Sleep(backoff)
		}
	}
	if err != nil { // 5. 所有重试都失败了，返回最后一次的错误
		return err
	}
	return nil
}

func NeedReDoForNetError(err error) bool {
	yes := err.Error() == io.EOF.Error() // 网络断开
	return yes
}
