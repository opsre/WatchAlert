package ctx

import (
	"context"
	"sync"
	"watchAlert/internal/cache"
	"watchAlert/internal/repo"
)

type Context struct {
	DB                 repo.InterEntryRepo
	Cache              cache.InterEntryCache
	Ctx                context.Context
	Mux                sync.RWMutex
	ConsumerContextMap map[string]context.CancelFunc
	CacheType          string
}

var (
	DB        repo.InterEntryRepo
	Cache     cache.InterEntryCache
	Ctx       context.Context
	CacheType string
)

func NewContext(ctx context.Context, db repo.InterEntryRepo, c cache.InterEntryCache, cacheType string) *Context {
	DB = db
	Cache = c
	Ctx = ctx
	CacheType = cacheType
	return &Context{
		DB:                 db,
		Cache:              c,
		Ctx:                ctx,
		ConsumerContextMap: make(map[string]context.CancelFunc),
		CacheType:          cacheType,
	}
}

func DO() *Context {
	return &Context{
		DB:        DB,
		Cache:     Cache,
		Ctx:       Ctx,
		CacheType: CacheType,
	}
}
