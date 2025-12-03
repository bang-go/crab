package crab

import (
	"context"
	"sync"

	"github.com/bang-go/crab/cmd"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/core/base/types"
	"github.com/bang-go/crab/core/pub/bag"
	"github.com/bang-go/crab/core/pub/graceful"
	"github.com/bang-go/crab/internal/log"
	"github.com/bang-go/crab/internal/vars"
	"github.com/bang-go/opt"
	"github.com/spf13/cobra"
)

type Handler struct {
	Pre   types.FuncErr //预加载
	Init  types.FuncErr //初始化
	Close types.FuncErr //关闭
}

type Worker interface {
	Use(...Handler) error
	RegisterCmd(...cmd.Cmder)
	Start() error
	Close()
	Done()
}

type artisanEntity struct {
	ctx         context.Context
	opt         *options
	preBagger   bag.Bagger
	initBagger  bag.Bagger
	closeBagger bag.Bagger
	done        chan struct{}
	Cmds        []cmd.Cmder
	commandPtr  *cobra.Command
	appName     string
}

var art *artisanEntity
var _ Worker = art
var m sync.RWMutex

// Build creates a new ant instance.
func Build(opts ...opt.Option[options]) {
	var err error
	o := &options{logOptions: logOptions{allowLogLevel: logx.LevelInfo, logEncodeType: logx.LogEncodeJson}, appName: vars.DefaultAppName.Load()}
	opt.Each(o, opts...)
	art = &artisanEntity{ctx: context.Background(),
		opt:         o,
		preBagger:   bag.NewBagger(),
		initBagger:  bag.NewBagger(),
		closeBagger: bag.NewBagger(),
		done:        make(chan struct{}, 1),
		commandPtr:  cmd.RootCmd,
		appName:     o.appName,
	}
	if err = art.init(); err != nil { //框架本身初始化加载
		panic(err)
	}
}

func defaultArtisan() *artisanEntity {
	if art == nil {
		Build()
	}
	return art
}

func Start() error {
	m.RLock()
	defer m.RUnlock()
	return defaultArtisan().Start()
}
func (a *artisanEntity) Start() error {
	go graceful.WatchSignal(a.done, a.closeBagger)
	if err := a.initBagger.Finish(); err != nil {
		return err
	}
	if len(a.Cmds) > 0 {
		for _, v := range a.Cmds {
			a.commandPtr.AddCommand(v.Cmd())
		}
		return a.commandPtr.Execute()
	}
	return nil
}

func (a *artisanEntity) init() error {
	//初始化日志客户端
	log.SetFrameLogger(a.opt.allowLogLevel, a.opt.logEncodeType)
	return nil
}

func Close() {
	m.RLock()
	defer m.RUnlock()
	defaultArtisan().Close()
}

// Close 停止
func (a *artisanEntity) Close() {
	//框架相关
	//应用相关
	if err := a.closeBagger.Finish(); err != nil {
		log.DefaultFrameLogger().Error(err.Error())
	}
}

func Use(Handlers ...Handler) error {
	m.RLock()
	defer m.RUnlock()
	return defaultArtisan().Use(Handlers...)
}

func (a *artisanEntity) Use(handlers ...Handler) error {
	for _, handler := range handlers {
		if handler.Pre != nil {
			a.preBagger.Register(handler.Pre)
			//直接运行pre
			err := handler.Pre()
			if err != nil {
				log.DefaultFrameLogger().Error(err.Error())
				return err
			}
		}
		if handler.Init != nil {
			a.initBagger.Register(handler.Init)
		}
		if handler.Close != nil {
			a.closeBagger.Register(handler.Close)
		}
	}

	return nil
}

func RegisterCmd(cmds ...cmd.Cmder) {
	m.RLock()
	defer m.RUnlock()
	defaultArtisan().RegisterCmd(cmds...)
}

func (a *artisanEntity) RegisterCmd(cmds ...cmd.Cmder) {
	a.Cmds = append(a.Cmds, cmds...)
}

func RegisterInitBagger(f ...types.FuncErr) {
	m.RLock()
	defer m.RUnlock()
	defaultArtisan().RegisterInitBagger(f...)
}

func (a *artisanEntity) RegisterInitBagger(f ...types.FuncErr) {
	a.initBagger.Register(f...)
}

func RegisterCloseBagger(f ...types.FuncErr) {
	m.RLock()
	defer m.RUnlock()
	defaultArtisan().RegisterCloseBagger(f...)
}

func (a *artisanEntity) RegisterCloseBagger(f ...types.FuncErr) {
	a.closeBagger.Register(f...)
}

func Done() {
	defaultArtisan().Done()
}

func (a *artisanEntity) Done() {
	a.done <- struct{}{}
}
