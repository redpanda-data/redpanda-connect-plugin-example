package condition

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/Jeffail/benthos/lib/cache"
	"github.com/Jeffail/benthos/lib/log"
	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/types"
)

//------------------------------------------------------------------------------

func init() {
	cache.RegisterPlugin(
		"lossy_memory",
		func() interface{} {
			conf := NewLossyMemoryConfig()
			return &conf
		},
		func(
			iconf interface{},
			mgr types.Manager,
			logger log.Modular,
			stats metrics.Type,
		) (types.Cache, error) {
			conf, ok := iconf.(*LossyMemoryConfig)
			if !ok {
				return nil, errors.New("failed to cast config")
			}
			return NewLossyMemory(*conf, mgr, logger, stats)
		},
	)
	cache.DocumentPlugin(
		"lossy_memory",
		`Attempts to stores cache items in memory but occasionally loses them,
oops! If the field `+"`max_size`"+` is greater than 0 then the total number of
keys will be capped at this value, items will be randomly removed whenever the
cap is reached.`,
		nil,
	)
}

//------------------------------------------------------------------------------

// LossyMemoryConfig is a configuration struct containing fields for the LossyMemory
// condition.
type LossyMemoryConfig struct {
	MaxSize int `json:"max_size" yaml:"max_size"`
}

// NewLossyMemoryConfig returns a LossyMemoryConfig with default values.
func NewLossyMemoryConfig() LossyMemoryConfig {
	return LossyMemoryConfig{
		MaxSize: 0,
	}
}

//------------------------------------------------------------------------------

// LossyMemory is a cache that keeps a lossy store of items in memory.
type LossyMemory struct {
	cap   int
	items map[string][]byte

	mKeys metrics.StatGauge

	sync.RWMutex
}

// NewLossyMemory returns a LossyMemory cache.
func NewLossyMemory(
	conf LossyMemoryConfig, mgr types.Manager, log log.Modular, stats metrics.Type,
) (types.Cache, error) {
	return &LossyMemory{
		cap:   conf.MaxSize,
		items: map[string][]byte{},
		mKeys: stats.GetGauge("keys"),
	}, nil
}

//------------------------------------------------------------------------------

// Get attempts to locate and return a cached value by its key, returns an error
// if the key does not exist.
func (m *LossyMemory) Get(key string) ([]byte, error) {
	m.RLock()
	k, exists := m.items[key]
	m.RUnlock()
	if !exists {
		return nil, types.ErrKeyNotFound
	}
	return k, nil
}

// Set attempts to set the value of a key.
func (m *LossyMemory) Set(key string, value []byte) error {
	if rand.Int()%7 == 0 {
		// Ooops!
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
	m.mKeys.Set(int64(len(m.items)))
	m.Unlock()
	return nil
}

// SetMulti attempts to set the value of multiple keys, returns an error if any
// keys fail.
func (m *LossyMemory) SetMulti(items map[string][]byte) error {
	m.Lock()
	for k, v := range items {
		if rand.Int()%7 == 0 {
			// Ooops!
			continue
		}
		m.items[k] = v
	}
	if m.cap > 0 && len(m.items) >= m.cap {
		for k := range m.items {
			if len(m.items) < m.cap && rand.Int()%5 == 0 {
				break
			}
			delete(m.items, k)
		}
	}
	m.mKeys.Set(int64(len(m.items)))
	m.Unlock()
	return nil
}

// Add attempts to set the value of a key only if the key does not already exist
// and returns an error if the key already exists.
func (m *LossyMemory) Add(key string, value []byte) error {
	m.Lock()
	if _, exists := m.items[key]; exists {
		m.Unlock()
		return types.ErrKeyAlreadyExists
	}
	if rand.Int()%7 == 0 {
		// Ooops!
		return nil
	}
	m.items[key] = value
	m.mKeys.Set(int64(len(m.items)))
	m.Unlock()
	return nil
}

// Delete attempts to remove a key.
func (m *LossyMemory) Delete(key string) error {
	m.Lock()
	delete(m.items, key)
	m.mKeys.Set(int64(len(m.items)))
	m.Unlock()
	return nil
}

// CloseAsync shuts down the cache.
func (m *LossyMemory) CloseAsync() {
}

// WaitForClose blocks until the cache has closed down.
func (m *LossyMemory) WaitForClose(timeout time.Duration) error {
	return nil
}

//------------------------------------------------------------------------------
