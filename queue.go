package funcs

import (
	"slices"
	"sync"
)

type QueueWithLength[T comparable] struct {
	lock   sync.Mutex
	Length int
	queue  []T
}

func (fs *QueueWithLength[T]) Push(inqueue T) (pops []T) {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	if slices.Contains(fs.queue, inqueue) { // 已存在，不重复添加
		return
	}
	fs.queue = append(fs.queue, inqueue)
	fileLength := len(fs.queue)
	if fs.Length > 0 && fileLength > fs.Length {
		pops = fs.queue[0 : fileLength-fs.Length]
		fs.queue = fs.queue[fileLength-fs.Length:]
	}

	return pops
}
