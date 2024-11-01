package graceful

import (
	"github.com/bang-go/crab/core/base/types"
	"github.com/bang-go/crab/core/pub/bag"
	"github.com/bang-go/crab/internal/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	MaxWaitTime time.Duration = 60 * time.Second
)

var shutdownBag = bag.NewBagger()

func Register(f ...types.FuncErr) {
	shutdownBag.Register(f...)
}

func WatchSignal(done chan struct{}, extBagger ...bag.Bagger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	select {
	case <-done:
		break
	case s := <-sigChan:
		log.DefaultFrameLogger().Warn("received signal", "sig", s.String())
		break
	}
	bagger := append(extBagger, shutdownBag)
	gracefulShutdown(sigChan, bagger...)
}

func gracefulShutdown(sig chan os.Signal, bagger ...bag.Bagger) {
	signal.Stop(sig)
	ch := make(chan struct{}, 1)
	go func() {
		for _, b := range bagger {
			_ = b.Finish()
		}
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		log.DefaultFrameLogger().Warn("graceful stop finish")
		break
	case <-time.After(MaxWaitTime):
		log.DefaultFrameLogger().Warn("The maximum wait time exceeded")
		break
	}
	//_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM) //该函数不支持windows
	pro, _ := os.FindProcess(syscall.Getpid()) //为了支持windows编译
	_ = pro.Kill()
}
