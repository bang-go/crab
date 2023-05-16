package health

import "sync/atomic"

type ProbeHealth interface {
	SetReady(status bool)
	SetLive(status bool)
	IsReady() bool
	IsLive() bool
}

var Probe ProbeHealth

func init() {
	Probe = NewHealth()
	Probe.SetLive(true)
}

func NewHealth() ProbeHealth {
	return &healthEntity{}
}

type healthEntity struct {
	readyStatus atomic.Bool
	liveStatus  atomic.Bool
}

func (h *healthEntity) SetReady(status bool) {
	h.readyStatus.Store(status)
}
func (h *healthEntity) SetLive(status bool) {
	h.liveStatus.Store(status)
}

func (h *healthEntity) IsReady() bool {
	return h.readyStatus.Load()
}
func (h *healthEntity) IsLive() bool {
	return h.liveStatus.Load()
}
