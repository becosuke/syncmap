package syncmap

import (
	"sync"
)

type Syncmap interface {
	Create(key, value any) error
	Read(key any) (any, error)
	Update(key, value any) error
	Delete(key any) error
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

func (impl *syncmapImpl) Read(key any) (any, error) {
	if key == nil {
		return nil, ErrInvalidArgument
	}
	value, ok := impl.syncmap.Load(key)
	if !ok {
		return nil, ErrNotFound
	}
	return value, nil
}

func (impl *syncmapImpl) Update(key, value any) error {
	if key == nil || value == nil {
		return ErrInvalidArgument
	}
	_, ok := impl.syncmap.Load(key)
	if !ok {
		return ErrNotFound
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
