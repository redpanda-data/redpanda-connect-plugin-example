package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/benthosdev/benthos/v4/public/service"
)

func init() {
	confSpec := service.NewConfigSpec().
		Summary("Creates a terrible cache with a fixed capacity.").
		Field(service.NewIntField("capacity").Default(100))

	err := service.RegisterCache(
		"lossy_memory", confSpec,
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.Cache, error) {
			capacity, err := conf.FieldInt("capacity")
			if err != nil {
				return nil, err
			}
			return newLossyMemory(capacity, mgr.Metrics())
		})
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type lossyMemory struct {
	cap   int
	items map[string][]byte

	mDropped *service.MetricCounter

	sync.RWMutex
}

func newLossyMemory(cap int, metrics *service.Metrics) (service.Cache, error) {
	return &lossyMemory{
		cap:      cap,
		items:    map[string][]byte{},
		mDropped: metrics.NewCounter("dropped.just.cus"),
	}, nil
}

func (m *lossyMemory) Get(ctx context.Context, key string) ([]byte, error) {
	m.RLock()
	k, exists := m.items[key]
	m.RUnlock()
	if !exists {
		return nil, service.ErrKeyNotFound
	}
	return k, nil
}

func (m *lossyMemory) Set(ctx context.Context, key string, value []byte, ttl *time.Duration) error {
	if rand.Int()%7 == 0 {
		// Ooops!
		m.mDropped.Incr(1)
		return nil
	}
	m.Lock()
	m.items[key] = value
	if m.cap > 0 && len(m.items) >= m.cap {
		for k := range m.items {
			if len(m.items) < m.cap && rand.Int()%5 == 0 {
				break
			}
			delete(m.items, k)
		}
	}
	m.Unlock()
	return nil
}

func (m *lossyMemory) Add(ctx context.Context, key string, value []byte, ttl *time.Duration) error {
	m.Lock()
	if _, exists := m.items[key]; exists {
		m.Unlock()
		return service.ErrKeyAlreadyExists
	}
	if rand.Int()%7 == 0 {
		// Ooops!
		m.mDropped.Incr(1)
		return nil
	}
	m.items[key] = value
	m.Unlock()
	return nil
}

func (m *lossyMemory) Delete(ctx context.Context, key string) error {
	m.Lock()
	delete(m.items, key)
	m.Unlock()
	return nil
}

func (m *lossyMemory) Close(ctx context.Context) error {
	return nil
}
