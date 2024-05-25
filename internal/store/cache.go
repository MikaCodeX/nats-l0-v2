package store

import (
	"sync"
)

type Cache struct {
	mx sync.RWMutex
	m  map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		m: make(map[string][]byte),
	}
}

func (c *Cache) Set(key string, value []byte) {
	c.mx.Lock()
	c.m[key] = value
	c.mx.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mx.RLock()
	val, ok := c.m[key]
	c.mx.RUnlock()
	return val, ok
}
