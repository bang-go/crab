package throttle_control

import (
	sentinelApi "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/internal/log"
	"github.com/bang-go/opt"
)

type limiter struct{}

func Limiter() ThrottlerLimiter {
	return &limiter{}
}

func (l *limiter) Build(opts ...opt.Option[options]) error {
	o := defaultOptions()
	opt.Each(o, opts...)
	conf := config.NewDefaultConfig() //todo: 增加更多options
	// default, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	logging.ResetGlobalLoggerLevel(o.logLevel)
	err := sentinelApi.InitWithConfig(conf)
	if err != nil {
		return err
	}
	return nil
}

// Rule 流量控制规则
func (l *limiter) Rule(rules []*flow.Rule) error {
	_, err := flow.LoadRules(rules)
	return err
}

func (l *limiter) Guard(resource string, pass FuncWithErr, reject Func, opts ...sentinelApi.EntryOption) bool {
	e, b := sentinelApi.Entry(resource, opts...)
	if b != nil {
		// Blocked. We could get the block reason from the BlockError.
		//log.Printf("sentinel throttle reject: %v", b.BlockMsg())
		log.FrameLogger.Warn("sentinel throttle reject: ", logx.String("msg", b.BlockMsg()))
		reject()
		return false
	} else {
		// Passed, wrap the logic here.
		_ = pass()
		e.Exit()
		return true
	}
}
