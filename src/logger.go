package log

import (
	"context"
	"log/slog"
	"sync/atomic"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

var (
	defaultLogger atomic.Pointer[Logger]
	_             Logger          = (*logger)(nil)
	_builder      InternalBuilder = &builder{}
)

func init() {
	SetDefault(GetBuilder().Develop())
}

type (
	Leveler          = slog.Leveler
	LevelVar         = slog.LevelVar
	Level            = slog.Level
	CallerFormatter  func(file string, line int) (repFile, repLine string)
	MessageFormatter func(message string) string
)

// GetBuilder 获取用于构建 Logger 的 Builder
func GetBuilder() InternalBuilder {
	return _builder
}

// InternalBuilder 是一个内部的日志记录器构建器，它提供了一些方法用于构建不同环境下的日志记录器
type InternalBuilder interface {
	Builder

	// DevelopOnGoland 构建一个适用于在 Goland 中开发时使用的日志记录器
	DevelopOnGoland() Logger
}

// Builder 是一个日志记录器构建器，它提供了一些方法用于构建不同环境下的日志记录器
type Builder interface {
	// Build 构建一个默认的日志记录器
	Build() Logger

	// FromConfiguration 以指定的选项构建一个日志记录器
	FromConfiguration(configuration LoggerConfiguration) Logger

	// FromConfigurators 以指定的配置器构建一个日志记录器
	FromConfigurators(configurators ...LoggerConfigurator) Logger

	// BuildWith 以指定的配置器构建一个日志记录器
	BuildWith(configuration LoggerConfiguration, configurators ...LoggerConfigurator) Logger

	// Develop 构建一个适用于开发环境的日志记录器
	Develop() Logger

	// Production 构建一个适用于生产环境的日志记录器
	Production() Logger

	// Test 构建一个适用于测试环境的日志记录器
	Test() Logger

	// Silent 构建一个静默的日志记录器，它不会输出任何日志
	Silent() Logger

	// Multi 构建一个多重日志记录器，它会将多个日志记录器组合在一起
	Multi(loggers ...Logger) Logger
}

type builder struct{}

func (b *builder) Multi(loggers ...Logger) Logger {
	handlers := make([]slog.Handler, 0, len(loggers))
	for _, l := range loggers {
		handlers = append(handlers, l.Handler())
	}
	return &logger{
		slog: slog.New(newMultiHandler(handlers...)),
	}
}

func (b *builder) Silent() Logger {
	return &logger{
		slog: slog.New(newSilentHandler()),
	}
}

func (b *builder) FromConfiguration(configuration LoggerConfiguration) Logger {
	return &logger{
		slog: slog.New(newHandler(configuration)),
	}
}

func (b *builder) FromConfigurators(configurators ...LoggerConfigurator) Logger {
	c := GetConfigBuilder().Build()
	for _, f := range configurators {
		f.Configure(c)
	}
	return &logger{
		slog: slog.New(newHandler(c.(LoggerOptionsFetcher))),
	}
}

func (b *builder) BuildWith(configuration LoggerConfiguration, configurators ...LoggerConfigurator) Logger {
	for _, f := range configurators {
		f.Configure(configuration)
	}
	return &logger{
		slog: slog.New(newHandler(configuration)),
	}
}

func (b *builder) Build() Logger {
	return &logger{
		slog: slog.New(newHandler(GetConfigBuilder().Build().(LoggerOptionsFetcher))),
	}
}

func (b *builder) Develop() Logger {
	return &logger{
		slog: slog.New(newHandler(GetConfigBuilder().Develop().(LoggerOptionsFetcher))),
	}
}

func (b *builder) DevelopOnGoland() Logger {
	return &logger{
		slog: slog.New(newHandler(GetConfigBuilder().DevelopOnGoland().(LoggerOptionsFetcher))),
	}
}

func (b *builder) Production() Logger {
	return &logger{
		slog: slog.New(newHandler(GetConfigBuilder().Production().(LoggerOptionsFetcher))),
	}
}

func (b *builder) Test() Logger {
	return &logger{
		slog: slog.New(newHandler(GetConfigBuilder().Test().(LoggerOptionsFetcher))),
	}
}

type Logger interface {
	// GetSLogger 获取 slog.Logger
	GetSLogger() *slog.Logger

	// Handler 获取 Handler
	Handler() Handler

	// With 为日志记录器创建一个新的实例
	With(args ...any) Logger

	// WithGroup 为日志记录器创建一个新的实例，并设置组名
	WithGroup(name string) Logger

	// Enabled 检查日志记录器是否启用了指定级别
	Enabled(ctx context.Context, level Level) bool

	// Log 在指定级别下记录一条消息
	Log(ctx context.Context, level Level, msg string, args ...any)

	// LogAttrs 在指定级别下记录一条消息，并附加属性
	LogAttrs(ctx context.Context, level Level, msg string, attrs ...Attr)

	// Debug 在 LevelDebug 级别下记录一条消息
	Debug(msg string, args ...any)

	// DebugContext 在 LevelDebug 级别下记录一条消息，并附加上下文
	DebugContext(ctx context.Context, msg string, args ...any)

	// Info 在 LevelInfo 级别下记录一条消息
	Info(msg string, args ...any)

	// InfoContext 在 LevelInfo 级别下记录一条消息，并附加上下文
	InfoContext(ctx context.Context, msg string, args ...any)

	// Warn 在 LevelWarn 级别下记录一条消息
	Warn(msg string, args ...any)

	// WarnContext 在 LevelWarn 级别下记录一条消息，并附加上下文
	WarnContext(ctx context.Context, msg string, args ...any)

	// Error 在 LevelError 级别下记录一条消息
	Error(msg string, args ...any)

	// ErrorContext 在 LevelError 级别下记录一条消息，并附加上下文
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type logger struct {
	slog *slog.Logger
}

func (l *logger) Enabled(ctx context.Context, level Level) bool {
	return l.slog.Enabled(ctx, level)
}

func (l *logger) Log(ctx context.Context, level Level, msg string, args ...any) {
	l.slog.Log(ctx, level, msg, args...)
}

func (l *logger) LogAttrs(ctx context.Context, level Level, msg string, attrs ...Attr) {
	l.slog.LogAttrs(ctx, level, msg, attrs...)
}

func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

func (l *logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.slog.DebugContext(ctx, msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

func (l *logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.slog.InfoContext(ctx, msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

func (l *logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.slog.WarnContext(ctx, msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

func (l *logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.slog.ErrorContext(ctx, msg, args...)
}

func (l *logger) GetSLogger() *slog.Logger {
	return l.slog
}

func (l *logger) Handler() Handler {
	return l.slog.Handler()
}

func (l *logger) clone() *logger {
	c := *l
	return &c
}

func (l *logger) With(args ...any) Logger {
	c := l.clone()
	c.slog = l.slog.With(args...)
	return c
}

func (l *logger) WithGroup(name string) Logger {
	c := l.clone()
	c.slog = l.slog.WithGroup(name)
	return c
}

// SetDefault 设置默认日志记录器
func SetDefault(logger Logger) {
	slog.Default().Handler()
	defaultLogger.Store(&logger)
}

// GetDefault 获取默认日志记录器
func GetDefault() Logger {
	l := defaultLogger.Load()
	return *l
}

// Debug 使用全局日志记录器在 LevelDebug 级别下记录一条消息
func Debug(msg string, args ...any) {
	l := *defaultLogger.Load()
	l.Debug(msg, args...)
}

// DebugContext 使用全局日志记录器在 LevelDebug 级别下记录一条消息，并附加上下文
func DebugContext(ctx context.Context, msg string, args ...any) {
	l := *defaultLogger.Load()
	l.DebugContext(ctx, msg, args...)
}

// Info 使用全局日志记录器在 LevelInfo 级别下记录一条消息
func Info(msg string, args ...any) {
	l := *defaultLogger.Load()
	l.Info(msg, args...)
}

// InfoContext 使用全局日志记录器在 LevelInfo 级别下记录一条消息，并附加上下文
func InfoContext(ctx context.Context, msg string, args ...any) {
	l := *defaultLogger.Load()
	l.InfoContext(ctx, msg, args...)
}

// Warn 使用全局日志记录器在 LevelWarn 级别下记录一条消息
func Warn(msg string, args ...any) {
	l := *defaultLogger.Load()
	l.Warn(msg, args...)
}

// WarnContext 使用全局日志记录器在 LevelWarn 级别下记录一条消息，并附加上下文
func WarnContext(ctx context.Context, msg string, args ...any) {
	l := *defaultLogger.Load()
	l.WarnContext(ctx, msg, args...)
}

// Error 使用全局日志记录器在 LevelError 级别下记录一条消息
func Error(msg string, args ...any) {
	l := *defaultLogger.Load()
	l.Error(msg, args...)
}

// ErrorContext 使用全局日志记录器在 LevelError 级别下记录一条消息，并附加上下文
func ErrorContext(ctx context.Context, msg string, args ...any) {
	l := *defaultLogger.Load()
	l.ErrorContext(ctx, msg, args...)
}
