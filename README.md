# Go-Log 日志库

一个简洁高效的 Go 语言日志库，支持多环境配置、灵活扩展和多日志实例组合，适用于开发、测试和生产环境。

---

## 安装

```bash
go get -u github.com/kercylan98/go-log
```

---

## 快速开始

以下示例展示了如何使用 Go-Log 日志库：

```go
package main

import (
	"errors"
	"github.com/kercylan98/go-log/log"
	"os"
)

func main() {
	var loggerA, loggerB log.Logger
	var builder = log.GetBuilder()

	// 使用默认日志
	loggerA = log.GetDefault()

	// 切换日志环境
	loggerA = builder.Build()         // 默认环境
	loggerA = builder.Develop()       // 开发环境
	loggerA = builder.Test()          // 测试环境
	loggerA = builder.Production()    // 生产环境
	loggerA = builder.Silent()        // 静默模式
	loggerA = builder.DevelopOnGoland() // 开发工具适配

	// 配置日志实例
	loggerB = builder.FromConfiguration(
		log.GetConfigBuilder().Build().
			WithWriter(os.Stdout).         // 设置输出目标
			WithLeveler(log.LevelInfo),    // 设置日志级别
	)

	loggerB = builder.FromConfigurators(log.LoggerConfiguratorFn(func(config log.LoggerConfiguration) {
		config.WithEnableColor(true).    // 启用彩色输出
			WithLevelStr(log.LevelDebug, "debug") // 自定义日志级别名称
	}))

	// 组合多个日志实例
	loggerC := builder.Multi(loggerA, loggerB)

	// 记录日志
	loggerC.Error("Example", "err", errors.New("example error"))
	loggerC.Error("Example", log.Err(errors.New("example error")))
}
```

---

## 许可证

MIT License. 详细信息请查看 [LICENSE](./LICENSE) 文件。

---

## 贡献

欢迎贡献代码或提出问题！如需贡献，请提交 Pull Request 或 Issue。
