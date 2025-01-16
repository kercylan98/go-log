package log

// Provider 是日志记录器提供器
type Provider interface {
	// Provide 提供日志记录器
	Provide() Logger
}

// ProviderFn 是一个函数类型的 Provider
type ProviderFn func() Logger

func (f ProviderFn) Provide() Logger {
	return f()
}
