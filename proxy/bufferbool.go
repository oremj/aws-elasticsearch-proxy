package proxy

import "sync"

type BufferPool struct {
	pool sync.Pool
}

func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 32*1024)
			},
		},
	}
}

func (b *BufferPool) Get() []byte {
	return b.pool.Get().([]byte)
}

func (b *BufferPool) Put(x []byte) {
	b.pool.Put(x)
}
