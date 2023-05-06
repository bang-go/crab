package crab

import (
	"github.com/bang-go/crab/cmd"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/internal/log"
	"github.com/bang-go/opt"
	"github.com/spf13/cobra"
	"sync"
)

type HandlerFunc func() error
type Handler struct {
	Name  string      //名称
	Init  HandlerFunc //初始化
	Close HandlerFunc //结束
}
type Worker interface {
	Exec([]Handler) error
	Use([]Handler)
	AddCmd(...cmd.Cmder)
	Start() error
	Exit() error
}

type artisan struct {
	opt          *options
	ExecHandlers []Handler
	UseHandlers  []Handler
	Cmds         []cmd.Cmder
	commandPtr   *cobra.Command
}
type logOptions struct {
	allowLogLevel logx.Level //允许的log level -1:Debug info:0 1:warn 2:error 3:dpanic 4 panic 5 fatal
	logEncoding   string     //日志编码 取值：json,console
}
type options struct {
	logOptions
}

var art *artisan
var m sync.RWMutex

// Build creates a new ant instance.
func Build(opts ...opt.Option[options]) {
	var err error
	o := &options{logOptions{allowLogLevel: logx.InfoLevel, logEncoding: logx.EncodeJson}}
	opt.Each(o, opts...)
	art = &artisan{opt: o, ExecHandlers: []Handler{}, UseHandlers: []Handler{}, commandPtr: cmd.RootCmd}
	if err = art.init(); err != nil { //框架预加载
		panic(err)
	}
}
func defaultArtisan() *artisan {
	return art
}

func Start() error {
	return defaultArtisan().start()
}
func (a *artisan) start() error {
	for _, h := range a.UseHandlers {
		if err := initHandler(h); err != nil {
			return err
		}
	}
	if len(a.Cmds) > 0 {
		for _, v := range a.Cmds {
			a.commandPtr.AddCommand(v.Cmd())
		}
		return a.commandPtr.Execute()
	}
	return nil
}

func (a *artisan) init() error {
	//初始化日志客户端
	log.InitLog(a.opt.allowLogLevel, a.opt.logEncoding)
	return nil
}

func Exit() error {
	return defaultArtisan().exit()
}

// Exit 停止
func (a *artisan) exit() error {
	//框架相关
	_ = log.FrameLogger.Sync()
	//应用相关
	exitHandle := append(a.ExecHandlers, a.UseHandlers...)
	if len(exitHandle) > 0 {
		for _, v := range exitHandle {
			if err := closeHandler(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func Exec(Handlers []Handler) error {
	m.RLock()
	defer m.RUnlock()
	return defaultArtisan().exec(Handlers)
}

// Exec 立刻会调用初始化函数，适用于需要立即执行Init，否则使用Use，按顺序加载
func (a *artisan) exec(Handlers []Handler) error {
	a.ExecHandlers = append(a.ExecHandlers, Handlers...)
	for _, h := range Handlers {
		if err := initHandler(h); err != nil {
			return err
		}
	}
	return nil
}

func Use(Handlers []Handler) {
	m.RLock()
	defer m.RUnlock()
	defaultArtisan().use(Handlers)
}
func (a *artisan) use(Handlers []Handler) {
	a.UseHandlers = append(a.UseHandlers, Handlers...)
}

func initHandler(h Handler) error {
	if h.Init == nil {
		return nil
	}
	if err := h.Init(); err != nil {
		log.FrameLogger.Error("init failed", logx.String("name", h.Name), logx.String("err", err.Error()))
		return err
	}
	log.FrameLogger.Info("init successful", logx.String("name", h.Name))
	return nil
}

func closeHandler(h Handler) error {
	if h.Close == nil {
		return nil
	}
	if err := h.Close(); err != nil {
		log.FrameLogger.Error("close failed", logx.String("name", h.Name), logx.String("err", err.Error()))
		return err
	}
	log.FrameLogger.Info("close successful", logx.String("name", h.Name))
	return nil
}

func AddCmd(cmds ...cmd.Cmder) {
	m.RLock()
	defer m.RUnlock()
	defaultArtisan().addCmd(cmds...)
}
func (a *artisan) addCmd(cmds ...cmd.Cmder) {
	a.Cmds = append(a.Cmds, cmds...)
}

func WithLogAllowLevel(logLevel logx.Level) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.allowLogLevel = logLevel
	})
}

func WithLogEncoding(logEncoding string) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.logEncoding = logEncoding
	})
}
