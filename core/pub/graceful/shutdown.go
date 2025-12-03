package graceful

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bang-go/crab/core/base/types"
	"github.com/bang-go/crab/core/pub/bag"
	"github.com/bang-go/crab/internal/log"
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
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-done:
		log.DefaultFrameLogger().Info("application done signal received")
	case s := <-sigChan:
		log.DefaultFrameLogger().Warn("received signal", "sig", s.String())
	}
	gracefulShutdown(sigChan, append(extBagger, shutdownBag)...)
}

func gracefulShutdown(sig chan os.Signal, bagger ...bag.Bagger) {
	signal.Stop(sig)
	ch := make(chan struct{}, 1)
	go func() {
		for _, b := range bagger {
			if err := b.Finish(); err != nil {
				log.DefaultFrameLogger().Error("error during shutdown", "error", err)
			}
		}
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		log.DefaultFrameLogger().Info("graceful shutdown completed")
	case <-time.After(MaxWaitTime):
		log.DefaultFrameLogger().Warn("graceful shutdown timeout exceeded")
	}
	pro, _ := os.FindProcess(syscall.Getpid())
	_ = pro.Kill()
}
