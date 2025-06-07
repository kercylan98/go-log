package main

import (
	"errors"
	"github.com/kercylan98/go-log/log"
	"os"
)

func main() {
	// 创建基本的日志记录器
	var builder = log.GetBuilder()
	baseLogger := builder.Develop()

	// 创建命名的日志记录器
	userLogger := baseLogger.WithGroup("user")
	authLogger := baseLogger.WithGroup("auth")
	dbLogger := baseLogger.WithGroup("database")

	// 使用命名的日志记录器
	userLogger.Info("用户登录成功", "userId", "12345")
	authLogger.Warn("检测到可疑的登录尝试", "ip", "192.168.1.100")
	dbLogger.Error("数据库连接失败", log.Err(errors.New("connection timeout")))

	// 演示嵌套命名
	userAuthLogger := userLogger.WithGroup("auth")
	userAuthLogger.Info("用户权限验证完成", "userId", "12345", "role", "admin")

	// 演示与 WithGroup 结合使用
	userGroupLogger := userLogger.WithGroup("session")
	userGroupLogger.Info("会话已创建",
		"sessionId", "abc123",
		"duration", "30m",
	)

	// 演示完整的日志配置示例
	configuredLogger := builder.FromConfiguration(
		log.GetConfigBuilder().Build().
			WithWriter(os.Stdout).
			WithLeveler(log.LevelInfo),
	)

	// 使用配置的日志记录器创建命名实例
	apiLogger := configuredLogger.WithGroup("api")
	apiLogger.Info("API服务启动成功", "port", 8080)

	// 演示多个日志记录器
	multiLogger := builder.Multi(
		apiLogger.WithGroup("service-a"),
		apiLogger.WithGroup("service-b"),
	)
	multiLogger.Error("服务间通信失败", log.Err(errors.New("timeout")))
}
