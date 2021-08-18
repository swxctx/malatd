package td

import (
	"github.com/swxctx/xlog"
)

const (
	TRCE    = "trace"
	DBUG    = "debug"
	INFO    = "info"
	WARN    = "warn"
	ERRO    = "error"
	FTAL    = "fatal"
	DISABLE = "disable"
)

// SetLevel 设置输出级别
func SetLevel(levelName string) {
	xlog.SetLevel(levelName)
}

// Fatalf
func Fatalf(format string, args ...interface{}) {
	xlog.Fatalf(format, args...)
}

// Errorf
func Errorf(format string, args ...interface{}) {
	xlog.Errorf(format, args...)
}

// Warnf
func Warnf(format string, args ...interface{}) {
	xlog.Warnf(format, args...)
}

// Infof
func Infof(format string, args ...interface{}) {
	xlog.Infof(format, args...)
}

// Debugf
func Debugf(format string, args ...interface{}) {
	xlog.Debugf(format, args...)
}

// Tracef
func Tracef(format string, args ...interface{}) {
	xlog.Tracef(format, args...)
}

// runLogPlugin
func runLogPlugin(ctx *Context) {
	Infof("From %s, %s", ctx.Request.RemoteAddr, ctx.Request.RequestURI)
	ctx.Next()
}
