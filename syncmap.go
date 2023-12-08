package syncmap

import (
	"sync"
)

type Syncmap interface {
	Create(key, value any) error
	Update(key, value any) error
	Delete(key any) error

	Get(key any) (any, error)
	GetMulti(keys []any) (map[any]any, error)
}

func NewSyncmap() Syncmap {
	return &syncmapImpl{
		syncmap: &sync.Map{},
	}
}

type syncmapImpl struct {
	syncmap *sync.Map
}

func (impl *syncmapImpl) Create(key, value any) error {
	if key == nil || value == nil {
		return ErrInvalidArgument
	}
	_, loaded := impl.syncmap.LoadOrStore(key, value)
	if loaded {
		return ErrAlreadyExists
	}
	return nil
}

func (impl *syncmapImpl) Update(key, value any) error {
	if key == nil || value == nil {
		return ErrInvalidArgument
	}
	impl.syncmap.Store(key, value)
	return nil
}

func (impl *syncmapImpl) Delete(key any) error {
	if key == nil {
		return ErrInvalidArgument
	}
	impl.syncmap.Delete(key)
	return nil
}

func (impl *syncmapImpl) Get(key any) (any, error) {
	if key == nil {
		return nil, ErrInvalidArgument
	}
	value, ok := impl.syncmap.Load(key)
	if !ok {
		return nil, ErrNotFound
	}
	return value, nil
}

func (impl *syncmapImpl) GetMulti(keys []any) (map[any]any, error) {
	if keys == nil {
		return nil, ErrInvalidArgument
	}
	res := make(map[any]any)
	for _, key := range keys {
		value, err := impl.Get(key)
		if err == nil {
			res[key] = value
		}
	}
	return res, nil
}
