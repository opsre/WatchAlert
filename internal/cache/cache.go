package cache

import (
	"time"
)

type Cache interface {
	SetKey(key, value string, expiration time.Duration)
	GetKey(key string) (string, error)
	DeleteKey(key string)
	SetHash(key, field, value string)
	SetHashAny(key, field string, value any)
	DeleteHash(key, field string)
	GetHash(key, field string) (string, error)
	GetHashAll(key string) (map[string]string, error)
}
