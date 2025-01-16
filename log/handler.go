package log

import (
	"context"
	"encoding"
	"fmt"
	"github.com/fatih/color"
	jsonIter "github.com/json-iterator/go"
	"github.com/kercylan98/go-log/log/internal/colorbuilder"
	"github.com/kercylan98/go-log/log/internal/convert"
	"log/slog"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

var (
	_ Handler = (*handler)(nil)
)

func newHandler(options LoggerOptionsFetcher) Handler {
	return &handler{
		options: options,
	}
}

// Handler 是基于 slog.Handler 的日志处理器
type Handler interface {
	slog.Handler // Handler 是 slog.Handler 的扩展
}

type handler struct {
	options       LoggerOptionsFetcher // options 是 Handler 的原始配置
	handleOptions LoggerOptionsFetcher // handleOptions 是 Handler 的运行时配置，它会在 Handle 时被复制

	attrs []slog.Attr
	group string
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.options.FetchLeveler().Level()
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	h.handleOptions = h.options.FetchCopy()
	if !h.Enabled(ctx, h.handleOptions.FetchLeveler().Level()) {
		return nil
	}

	var builder = colorbuilder.NewBuilder()
	defer builder.Reset()

	h.formatTime(ctx, record, builder)
	h.formatLevel(ctx, record, builder)
	h.formatCaller(ctx, record, builder)
	h.formatMessage(ctx, record, builder)

	// fixed attrs
	num := record.NumAttrs()
	fixedNum := len(h.attrs)
	for i, attr := range h.attrs {
		h.formatAttr(ctx, h.group, record.Level, attr, builder, num+fixedNum == i+1)
	}

	idx := 0
	record.Attrs(func(attr slog.Attr) bool {
		idx++
		h.formatAttr(ctx, h.group, record.Level, attr, builder, num == idx)
		return true
	})

	recordBytes, err := builder.Write('\n').Bytes()
	if err != nil {
		return err
	}

	_, err = h.handleOptions.FetchWriter().Write(recordBytes)
	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	n := h.clone()
	n.attrs = append(n.attrs, attrs...)
	return n
}

func (h *handler) WithGroup(name string) slog.Handler {
	n := h.clone()
	n.group = name
	return n
}

func (h *handler) clone() *handler {
	return &handler{
		options: h.options,
		attrs:   h.attrs,
		group:   h.group,
	}
}

func (h *handler) formatTime(ctx context.Context, record slog.Record, builder *colorbuilder.Builder) {
	h.loadAttrKey(builder, AttrKeyTime)
	h.loadColor(builder, ColorTypeTime).
		WriteString(record.Time.Format(h.handleOptions.FetchTimeLayout())).
		DisableColor().
		Write(' ')
}

func (h *handler) formatLevel(ctx context.Context, record slog.Record, builder *colorbuilder.Builder) {
	var colorType ColorType
	if h.handleOptions.FetchEnableColor() {
		switch record.Level {
		case slog.LevelDebug:
			colorType = ColorTypeDebugLevel
		case slog.LevelInfo:
			colorType = ColorTypeInfoLevel
		case slog.LevelWarn:
			colorType = ColorTypeWarnLevel
		case slog.LevelError:
			colorType = ColorTypeErrorLevel
		}
	}
	h.loadAttrKey(builder, AttrKeyLevel)
	h.loadColor(builder, colorType).
		WriteString(h.handleOptions.FetchLevelStr(record.Level)).
		DisableColor().
		Write(' ')
}

func (h *handler) formatCaller(ctx context.Context, record slog.Record, builder *colorbuilder.Builder) {
	if !h.handleOptions.FetchCaller() {
		return
	}
	pcs := make([]uintptr, 1)
	runtime.CallersFrames(pcs[:runtime.Callers(h.handleOptions.FetchCallerSkip(), pcs)])
	fs := runtime.CallersFrames(pcs)
	f, _ := fs.Next()
	if f.File == "" {
		return
	}

	var file, line string
	if callerFormatter := h.handleOptions.FetchCallerFormatter(); callerFormatter != nil {
		file, line = callerFormatter(f.File, f.Line)
	} else {
		file = filepath.Base(f.File)
		line = convert.IntToString(f.Line)
	}

	h.loadAttrKey(builder, AttrKeyCaller)
	h.loadColor(builder, ColorTypeCaller).
		WriteString(file).
		SetColor(h.handleOptions.FetchColorType(ColorTypeAttrDelimiter)).
		WriteString(":").
		SetColor(h.handleOptions.FetchColorType(ColorTypeAttrValue)).
		WriteString(line).
		DisableColor().
		Write(' ')
}

func (h *handler) formatMessage(ctx context.Context, record slog.Record, builder *colorbuilder.Builder) {
	if record.Message == "" {
		return
	}
	var msg = record.Message
	if messageFormatter := h.handleOptions.FetchMessageFormatter(); messageFormatter != nil {
		msg = messageFormatter(msg)
	}

	h.loadAttrKey(builder, AttrKeyMessage)
	h.loadColor(builder, ColorTypeMessage).
		WriteString(msg).
		DisableColor().
		Write(' ')
}

func (h *handler) formatAttr(ctx context.Context, group string, level slog.Level, attr slog.Attr, builder *colorbuilder.Builder, last bool) {
	var key = attr.Key
	if group != "" {
		key = group + "." + key
	}

	switch attr.Value.Kind() {
	case slog.KindGroup:
		groupAttr := attr.Value.Group()
		for _, a := range groupAttr {
			h.formatAttr(ctx, key, level, a, builder, last)
		}
		return
	default:
		h.loadColor(builder, ColorTypeAttrKey)
		switch v := attr.Value.Any().(type) {
		case stackError, stackErrorTracks:
			h.loadColor(builder, ColorTypeAttrErrorKey)
		case error:
			if h.handleOptions.FetchErrTrackLevel(level) && !h.handleOptions.FetchTrackBeautify() {
				pc := make([]uintptr, 10)
				n := runtime.Callers(h.handleOptions.FetchCallerSkip()+3, pc)
				frames := runtime.CallersFrames(pc[:n])
				var stacks = make(stackErrorTracks, 0, 10)
				for {
					frame, more := frames.Next()
					stacks = append(stacks, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
					if !more {
						break
					}
				}
				attr = slog.Group(attr.Key, slog.Any("info", stackError{v}), slog.Any("stack", stacks))
				h.formatAttr(ctx, group, level, attr, builder, false)
				return
			}
			h.loadColor(builder, ColorTypeAttrErrorKey)
		}
	}

	builder.
		WriteString(key).
		SetColor(h.handleOptions.FetchColorType(ColorTypeAttrDelimiter)).
		WriteString(h.handleOptions.FetchDelimiter())
	h.formatAttrValue(ctx, level, key, attr, builder, last)
}

func (h *handler) formatAttrValue(ctx context.Context, level slog.Level, fullKey string, attr slog.Attr, builder *colorbuilder.Builder, last bool) {
	h.loadColor(builder, ColorTypeAttrValue)
	defer builder.DisableColor()

	switch attr.Value.Kind() {
	case slog.KindString:
		builder.WriteString(strconv.Quote(attr.Value.String()))
	case slog.KindInt64:
		builder.WriteInt64(attr.Value.Int64())
	case slog.KindUint64:
		builder.WriteUint64(attr.Value.Uint64())
	case slog.KindFloat64:
		builder.WriteFloat64(attr.Value.Float64())
	case slog.KindBool:
		builder.WriteBool(attr.Value.Bool())
	case slog.KindDuration:
		builder.WriteString(strconv.Quote(attr.Value.Duration().String()))
	case slog.KindTime:
		builder.WriteString(strconv.Quote(attr.Value.Time().String()))
	default:
		switch v := attr.Value.Any().(type) {
		case stackError:
			h.loadColor(builder, ColorTypeAttrErrorKey)
			builder.WriteString(strconv.Quote(v.err.Error()))
		case stackErrorTracks:
			h.loadColor(builder, ColorTypeAttrErrorKey)
			builder.WriteString(strconv.Quote(fmt.Sprintf("%+v", attr.Value.Any())))
		case error:
			h.loadColor(builder, ColorTypeAttrErrorValue)
			builder.WriteString(strconv.Quote(v.Error()))

			if h.handleOptions.FetchErrTrackLevel(level) && h.handleOptions.FetchTrackBeautify() {
				pc := make([]uintptr, 10)
				n := runtime.Callers(h.handleOptions.FetchCallerSkip()+3, pc)
				frames := runtime.CallersFrames(pc[:n])
				if h.handleOptions.FetchTrackBeautify() {
					h.loadColor(builder, ColorTypeErrorTrackHeader).
						WriteSprintfToEnd("\tError Track: [%s] >> %s", fullKey, v.Error())
					h.loadColor(builder, ColorTypeErrorTrack)
					for {
						builder.WriteToEnd('\n')
						frame, more := frames.Next()
						builder.WriteToEnd('\t')
						builder.WriteStringToEnd(frame.File)
						builder.WriteToEnd(':')
						builder.WriteIntToEnd(frame.Line)
						builder.WriteToEnd(' ')
						builder.WriteStringToEnd(frame.Function)
						if !more {
							break
						}
					}
					builder.WriteToEnd('\n')
				}
			}
		case nil:
			builder.WriteString("<nil>")
		case encoding.TextMarshaler:
			data, err := v.MarshalText()
			if err != nil {
				break
			}
			builder.WriteString(strconv.Quote(string(data)))
		case []byte:
			builder.WriteString(strconv.Quote(*(*string)(unsafe.Pointer(&v))))
		case stack:
			if len(v) == 0 {
				builder.WriteString("<none>")
			} else {
				lines := strings.Split(string(v), "\n")
				builder.WriteString(fmt.Sprintf("lines(%d)", len(lines)))
				if h.handleOptions.FetchTrackBeautify() {
					for _, line := range lines {
						builder.WriteToEnd('\n')
						builder.WriteStringToEnd(line)
					}
					builder.WriteToEnd('\n')
				}
			}

		default:
			//builder.WriteString(strconv.Quote(fmt.Sprintf("%+v", attr.Values.Any())))
			jsonBytes, err := jsonIter.ConfigCompatibleWithStandardLibrary.Marshal(attr.Value.Any())
			if err != nil {
				jsonBytes = []byte("{}")
			}
			builder.WriteString(string(jsonBytes))
		}
	}

	if !last {
		builder.Write(' ')
	}
}

func (h *handler) loadColor(builder *colorbuilder.Builder, t ColorType) *colorbuilder.Builder {
	var c *color.Color
	if h.handleOptions.FetchEnableColor() {
		c = h.handleOptions.FetchColorType(t)
	}
	return builder.SetColor(c)
}

func (h *handler) loadAttrKey(builder *colorbuilder.Builder, key AttrKey) *colorbuilder.Builder {
	v, exist := h.handleOptions.FetchAttrKeys(key)
	if !exist {
		return builder
	}
	return builder.
		SetColor(h.handleOptions.FetchColorType(ColorTypeAttrKey)).
		WriteString(v).
		SetColor(h.handleOptions.FetchColorType(ColorTypeAttrDelimiter)).
		WriteString(h.handleOptions.FetchDelimiter())
}
