package log

import (
	"github.com/fatih/color"
	"github.com/kercylan98/go-log/src/internal/charproc"
	"github.com/kercylan98/go-log/src/internal/options"
	"io"
	"os"
	"sync"
	"time"
)

var (
	_               LoggerConfiguration  = (*loggerConfiguration)(nil)
	_optionsBuilder ConfigurationBuilder = new(configurationBuilder)
)

// GetConfigBuilder 获取一个 Logger 选项构建器
func GetConfigBuilder() ConfigurationBuilder {
	return _optionsBuilder
}

// ConfigurationBuilder 是 LoggerOptions 的构建器
type ConfigurationBuilder interface {
	// Build 构建一个默认的选项配置
	Build() LoggerConfiguration

	// Develop 构建一个适用于开发环境的选项配置
	Develop() LoggerConfiguration

	// DevelopOnGoland 构建一个适用于 Goland 开发环境的选项配置
	DevelopOnGoland() LoggerConfiguration

	// Test 构建一个适用于测试环境的选项配置
	Test() LoggerConfiguration

	// Production 构建一个适用于生产环境的选项配置
	Production() LoggerConfiguration
}

type configurationBuilder struct{}

func (o *configurationBuilder) Build() LoggerConfiguration {
	c := &loggerConfiguration{
		rw: new(sync.RWMutex),
	}
	c.LogicOptions = options.NewLogicOptions[LoggerOptionsFetcher, LoggerOptions](c, c)
	return c.
		WithWriter(os.Stdout).
		WithLeveler(LevelInfo).
		WithTimeLayout(time.DateTime).
		WithDelimiter("=").
		WithLevelStr(LevelDebug, "DBG").
		WithLevelStr(LevelInfo, "INF").
		WithLevelStr(LevelWarn, "WAR").
		WithLevelStr(LevelError, "ERR").
		WithColor(ColorTypeDebugLevel, color.FgHiCyan).
		WithColor(ColorTypeInfoLevel, color.FgHiGreen).
		WithColor(ColorTypeWarnLevel, color.FgHiYellow).
		WithColor(ColorTypeErrorLevel, color.FgHiRed).
		WithColor(ColorTypeMessage, color.FgHiBlack, color.Bold).
		WithColor(ColorTypeAttrDelimiter, color.FgHiBlack).
		WithColor(ColorTypeAttrKey, color.FgWhite).
		WithColor(ColorTypeAttrErrorKey, color.FgHiRed).
		WithColor(ColorTypeAttrErrorValue, color.FgHiRed).
		WithColor(ColorTypeErrorTrack, color.FgWhite).
		WithColor(ColorTypeErrorTrackHeader, color.FgYellow).
		WithCaller(true).
		WithCallerSkip(6).
		WithMessageFormatter(func(message string) string {
			return charproc.BigCamel(message)
		}).(LoggerConfiguration)
}

func (o *configurationBuilder) Develop() LoggerConfiguration {
	return o.Build().
		WithLeveler(LevelDebug).
		WithEnableColor(true).
		WithErrTrackLevel(LevelError).
		WithTrackBeautify(true).(LoggerConfiguration)
}

func (o *configurationBuilder) DevelopOnGoland() LoggerConfiguration {
	return o.Develop().
		WithLevelStr(LevelDebug, "DEBUG").
		WithLevelStr(LevelInfo, "INFO").
		WithLevelStr(LevelWarn, "WARN").
		WithLevelStr(LevelError, "ERROR").(LoggerConfiguration)
}

func (o *configurationBuilder) Test() LoggerConfiguration {
	return o.Develop().
		WithEnableColor(false).(LoggerConfiguration)
}

func (o *configurationBuilder) Production() LoggerConfiguration {
	return o.Build().
		WithLeveler(LevelInfo).
		WithEnableColor(false).(LoggerConfiguration)
}

// LoggerConfigurator 是 LoggerConfiguration 的配置接口，它允许结构化的配置 Logger
type LoggerConfigurator interface {
	Configure(config LoggerConfiguration) Logger
}

// LoggerConfiguratorFn 是 LoggerConfiguration 的配置接口，它允许通过函数式的方式配置 Logger
type LoggerConfiguratorFn func(config LoggerConfiguration) Logger

func (f LoggerConfiguratorFn) Configure(config LoggerConfiguration) Logger {
	return f(config)
}

// LoggerConfiguration 是 Logger 的配置接口，它支持运行时进行配置变更，并且是并发安全的
type LoggerConfiguration interface {
	LoggerOptions
	LoggerOptionsFetcher
}

type LoggerOptions interface {
	options.LogicOptions[LoggerOptionsFetcher, LoggerOptions]

	// WithTrackBeautify 设置错误追踪美化是否启用
	//   - 如果启用，那么当记录到 error 类型的日志时，将会得到易于阅读的错误追踪
	WithTrackBeautify(enable bool) LoggerOptions

	// WithErrTrackLevel 设置错误追踪级别，只有在指定的级别下才会记录错误追踪
	WithErrTrackLevel(levels ...Level) LoggerOptions

	// WithMessageFormatter 设置消息格式化器
	WithMessageFormatter(formatter MessageFormatter) LoggerOptions

	// WithCaller 设置是否显示调用者信息
	//  - 如果启用，那么将会显示调用者信息
	WithCaller(enable bool) LoggerOptions

	// WithCallerSkip 设置调用者跳过层数
	//  - 调用者跳过层数表示在获取调用者信息时，跳过的层数
	WithCallerSkip(skip int) LoggerOptions

	// WithCallerFormatter 设置调用者格式化器
	WithCallerFormatter(formatter CallerFormatter) LoggerOptions

	// WithLevelStr 设置日志级别所使用的字符串
	WithLevelStr(level Level, str string) LoggerOptions

	// WithDelimiter 设置分隔符
	WithDelimiter(delimiter string) LoggerOptions

	// WithAttrKey 设置属性键
	WithAttrKey(key AttrKey, str string) LoggerOptions

	// WithEnableColor 设置是否启用颜色
	WithEnableColor(enable bool) LoggerOptions

	// WithColor 设置日志颜色
	WithColor(colorType ColorType, attrs ...color.Attribute) LoggerOptions

	// WithTimeLayout 设置日志时间格式，如 "2006-01-02 15:04:05"
	WithTimeLayout(layout string) LoggerOptions

	// WithLeveler 设置日志级别
	WithLeveler(leveler Leveler) LoggerOptions

	// WithWriter 设置日志写入器
	WithWriter(writer io.Writer) LoggerOptions
}

type LoggerOptionsFetcher interface {
	// FetchTrackBeautify 获取错误追踪美化是否启用
	FetchTrackBeautify() bool

	// FetchLeveler 获取日志级别
	FetchLeveler() Leveler

	// FetchTimeLayout 获取日志时间格式
	FetchTimeLayout() string

	// FetchEnableColor 获取是否启用颜色
	FetchEnableColor() bool

	// FetchLevelStr 获取日志级别字符串
	FetchLevelStr(level Level) string

	// FetchCaller 获取是否显示调用者信息
	FetchCaller() bool

	// FetchCallerSkip 获取调用者跳过层数
	FetchCallerSkip() int

	// FetchCallerFormatter 获取调用者格式化器
	FetchCallerFormatter() CallerFormatter

	// FetchColorType 获取颜色类型
	FetchColorType(colorType ColorType) *color.Color

	// FetchMessageFormatter 获取消息格式化器
	FetchMessageFormatter() MessageFormatter

	// FetchErrTrackLevel 获取错误追踪级别
	FetchErrTrackLevel(level Level) bool

	// FetchDelimiter 获取分隔符
	FetchDelimiter() string

	// FetchAttrKeys 获取属性键
	FetchAttrKeys(key AttrKey) (v string, exist bool)

	// FetchCopy 获取一个副本
	FetchCopy() LoggerOptionsFetcher

	// FetchWriter 获取日志写入器
	FetchWriter() io.Writer
}

type loggerConfiguration struct {
	options.LogicOptions[LoggerOptionsFetcher, LoggerOptions]
	rw               *sync.RWMutex              // 读写锁
	leveler          Leveler                    // 日志级别
	timeLayout       string                     // 时间格式
	colorTypes       map[ColorType]*color.Color // 颜色类型
	enableColor      bool                       // 是否启用颜色
	attrKeys         map[AttrKey]string         // 属性键
	delimiter        string                     // 分隔符
	levelStr         map[Level]string           // 日志级别字符串
	caller           bool                       // 是否显示调用者
	callerSkip       int                        // 调用者跳过层数
	callerFormatter  CallerFormatter            // 调用者格式化
	messageFormatter MessageFormatter           // 消息格式化
	errTrackLevel    map[Level]struct{}         // 错误追踪级别
	trackBeautify    bool                       // 错误追踪美化
	writer           io.Writer                  // 日志写入器
}

func (h *loggerConfiguration) WithWriter(writer io.Writer) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.writer = writer
	})
}

func (h *loggerConfiguration) FetchWriter() (writer io.Writer) {
	h.fetch(func(config *loggerConfiguration) {
		writer = config.writer
	})
	return writer
}

func cloneMap[K comparable, V any](m map[K]V) map[K]V {
	clone := make(map[K]V)
	for k, v := range m {
		clone[k] = v
	}
	return clone

}

func (h *loggerConfiguration) FetchCopy() (fetcher LoggerOptionsFetcher) {
	h.fetch(func(config *loggerConfiguration) {
		clone := *config
		clone.rw = new(sync.RWMutex)

		clone.colorTypes = cloneMap(config.colorTypes)
		clone.attrKeys = cloneMap(config.attrKeys)
		clone.levelStr = cloneMap(config.levelStr)
		clone.errTrackLevel = cloneMap(config.errTrackLevel)

		fetcher = &clone
	})
	return

}

func (h *loggerConfiguration) FetchLeveler() (leveler Leveler) {
	h.fetch(func(config *loggerConfiguration) {
		leveler = config.leveler
	})
	return
}

func (h *loggerConfiguration) FetchTimeLayout() (layout string) {
	h.fetch(func(config *loggerConfiguration) {
		layout = config.timeLayout
	})
	return
}

func (h *loggerConfiguration) FetchEnableColor() (enable bool) {
	h.fetch(func(config *loggerConfiguration) {
		enable = config.enableColor
	})
	return
}

func (h *loggerConfiguration) FetchLevelStr(level Level) (str string) {
	h.fetch(func(config *loggerConfiguration) {
		str = config.levelStr[level]
	})
	return
}

func (h *loggerConfiguration) FetchCaller() (enable bool) {
	h.fetch(func(config *loggerConfiguration) {
		enable = config.caller
	})
	return
}

func (h *loggerConfiguration) FetchCallerSkip() (skip int) {
	h.fetch(func(config *loggerConfiguration) {
		skip = config.callerSkip
	})
	return
}

func (h *loggerConfiguration) FetchCallerFormatter() (formatter CallerFormatter) {
	h.fetch(func(config *loggerConfiguration) {
		formatter = config.callerFormatter
	})
	return
}

func (h *loggerConfiguration) FetchColorType(colorType ColorType) (c *color.Color) {
	h.fetch(func(config *loggerConfiguration) {
		c = config.colorTypes[colorType]
	})
	return
}

func (h *loggerConfiguration) FetchMessageFormatter() (formatter MessageFormatter) {
	h.fetch(func(config *loggerConfiguration) {
		formatter = config.messageFormatter
	})
	return
}

func (h *loggerConfiguration) FetchErrTrackLevel(level Level) (exist bool) {
	h.fetch(func(config *loggerConfiguration) {
		_, exist = config.errTrackLevel[level]
	})
	return
}

func (h *loggerConfiguration) FetchDelimiter() (delimiter string) {
	h.fetch(func(config *loggerConfiguration) {
		delimiter = config.delimiter
	})
	return
}

func (h *loggerConfiguration) FetchAttrKeys(key AttrKey) (v string, exist bool) {
	h.fetch(func(config *loggerConfiguration) {
		v, exist = config.attrKeys[key]
	})
	return
}

func (h *loggerConfiguration) update(logger func(config *loggerConfiguration)) *loggerConfiguration {
	h.rw.Lock()
	defer h.rw.Unlock()
	logger(h)
	return h
}

func (h *loggerConfiguration) fetch(logger func(config *loggerConfiguration)) {
	h.rw.RLock()
	defer h.rw.RUnlock()
	logger(h)
}

func (h *loggerConfiguration) WithTrackBeautify(enable bool) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.trackBeautify = enable
	})
}

func (h *loggerConfiguration) WithErrTrackLevel(levels ...Level) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		if config.errTrackLevel == nil {
			config.errTrackLevel = make(map[Level]struct{})
		}
		for _, level := range levels {
			config.errTrackLevel[level] = struct{}{}
		}
	})
}

func (h *loggerConfiguration) WithMessageFormatter(formatter MessageFormatter) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.messageFormatter = formatter
	})
}

func (h *loggerConfiguration) WithCaller(enable bool) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.caller = enable
	})
}

func (h *loggerConfiguration) WithCallerSkip(skip int) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.callerSkip = skip
	})
}

func (h *loggerConfiguration) WithCallerFormatter(formatter CallerFormatter) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.callerFormatter = formatter
	})
}

func (h *loggerConfiguration) WithLevelStr(level Level, str string) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		if config.levelStr == nil {
			config.levelStr = make(map[Level]string)
		}
		config.levelStr[level] = str
	})
}

func (h *loggerConfiguration) WithDelimiter(delimiter string) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.delimiter = delimiter
	})
}

func (h *loggerConfiguration) WithAttrKey(key AttrKey, str string) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		if config.attrKeys == nil {
			config.attrKeys = make(map[AttrKey]string)
		}
		config.attrKeys[key] = str
	})
}

func (h *loggerConfiguration) WithEnableColor(enable bool) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.enableColor = enable
	})
}

func (h *loggerConfiguration) WithColor(colorType ColorType, attrs ...color.Attribute) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		if config.colorTypes == nil {
			config.colorTypes = make(map[ColorType]*color.Color)
		}
		c := color.New(attrs...)
		c.EnableColor()
		config.colorTypes[colorType] = c
	})
}

func (h *loggerConfiguration) WithTimeLayout(layout string) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.timeLayout = layout
	})
}

func (h *loggerConfiguration) WithLeveler(leveler Leveler) LoggerOptions {
	return h.update(func(config *loggerConfiguration) {
		config.leveler = leveler
	})
}

func (h *loggerConfiguration) FetchTrackBeautify() bool {
	var result bool
	h.fetch(func(config *loggerConfiguration) {
		result = config.trackBeautify
	})
	return result
}
