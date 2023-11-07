package bag

import (
	"github.com/bang-go/crab/core/base/types"
	"sync"
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
	once sync.Once //确保同一个listen只会finish一次
}

func (b *BaggerEntity) Register(f ...types.FuncErr) {
	b.m.RLock()
	defer b.m.RUnlock()
	b.list = append(b.list, f...)
}

func (b *BaggerEntity) Finish() (err error) {
	b.once.Do(func() {
		for _, l := range b.list {
			err = l()
			if err != nil {
				return
			}
		}
	})
	return
}

func (b *BaggerEntity) Copy() Bagger {
	Basket := &BaggerEntity{}
	copy(Basket.list, b.list)
	return Basket
}
