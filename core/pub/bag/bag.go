package bag

import (
	"sync"

	"github.com/bang-go/crab/core/base/types"
)

type Bagger interface {
	Register(f ...types.FuncErr)
	Finish() error
}

func NewBagger() Bagger {
	return &BaggerEntity{}
}

type BaggerEntity struct {
	list []types.FuncErr
	m    sync.RWMutex
	once sync.Once
}

func (b *BaggerEntity) Register(f ...types.FuncErr) {
	b.m.Lock()
	defer b.m.Unlock()
	b.list = append(b.list, f...)
}

func (b *BaggerEntity) Finish() (err error) {
	b.once.Do(func() {
		b.m.RLock()
		list := make([]types.FuncErr, len(b.list))
		copy(list, b.list)
		b.m.RUnlock()
		for _, l := range list {
			err = l()
			if err != nil {
				return
			}
		}
	})
	return
}

func (b *BaggerEntity) Copy() Bagger {
	b.m.RLock()
	defer b.m.RUnlock()
	c := &BaggerEntity{}
	c.list = make([]types.FuncErr, len(b.list))
	copy(c.list, b.list)
	return c
}
