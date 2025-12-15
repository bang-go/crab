package logx

import (
	"context"
	"log/slog"
	"runtime"
)

// CallerHandler 包装 slog.Handler，支持自定义 caller skip 层数
type CallerHandler struct {
	handler slog.Handler
	skip    int
}

// NewCallerHandler 创建一个支持自定义 caller skip 的 Handler
func NewCallerHandler(handler slog.Handler, skip int) *CallerHandler {
	return &CallerHandler{
		handler: handler,
		skip:    skip,
	}
}

// GetSkip 获取当前的 skip 值
func (h *CallerHandler) GetSkip() int {
	return h.skip
}

// GetHandler 获取内部的 handler
func (h *CallerHandler) GetHandler() slog.Handler {
	return h.handler
}

func (h *CallerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CallerHandler) Handle(ctx context.Context, r slog.Record) error {
	// 获取正确的调用栈位置
	var pcs [1]uintptr
	// runtime.Callers 的 skip 参数：
	// 0 = Callers 本身
	// 1 = 当前函数 (Handle)
	// 2 = slog.Logger.log
	// 3 = slog.Logger.Info/Debug/...
	// 4+ = 用户代码
	// 所以我们需要 skip + 1（Handle 本身）
	runtime.Callers(h.skip+1, pcs[:])

	// 创建新的 Record 并设置正确的 PC
	r.PC = pcs[0]

	return h.handler.Handle(ctx, r)
}

func (h *CallerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CallerHandler{
		handler: h.handler.WithAttrs(attrs),
		skip:    h.skip,
	}
}

func (h *CallerHandler) WithGroup(name string) slog.Handler {
	return &CallerHandler{
		handler: h.handler.WithGroup(name),
		skip:    h.skip,
	}
}
