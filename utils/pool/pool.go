package pool

import "fmt"

type Pool struct {
	pos int
	buf []byte
}

const maxpoolsize = 500 * 1024

func (pool *Pool) Get(size int) []byte {
	if maxpoolsize-pool.pos < size {
		pool.pos = 0
		pool.buf = make([]byte, maxpoolsize)
	}
	if len(pool.buf) < pool.pos+size {
		fmt.Printf("err ")
	}
	b := pool.buf[pool.pos : pool.pos+size]
	pool.pos += size
	return b
}

func NewPool() *Pool {
	return &Pool{
		buf: make([]byte, maxpoolsize),
	}
}
