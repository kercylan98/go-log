package log

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"
)

// TestHandlerRaceCondition tests that there is no race condition when multiple goroutines
// access the handler's options concurrently.
func TestHandlerRaceCondition(t *testing.T) {
	// Create a logger configuration
	config := GetConfigBuilder().Build().
		WithLeveler(LevelDebug).
		WithTimeLayout("2006-01-02 15:04:05").
		WithEnableColor(true).
		WithCaller(true).
		WithCallerSkip(3).
		WithDelimiter(":").
		WithLevelStr(LevelDebug, "DEBUG").
		WithLevelStr(LevelInfo, "INFO").
		WithLevelStr(LevelWarn, "WARN").
		WithLevelStr(LevelError, "ERROR").
		WithAttrKey(AttrKeyTime, "time").
		WithAttrKey(AttrKeyLevel, "level").
		WithAttrKey(AttrKeyCaller, "caller").
		WithAttrKey(AttrKeyMessage, "msg").
		WithErrTrackLevel(LevelError).
		WithTrackBeautify(true).(LoggerOptionsFetcher)

	// Create a handler
	handler := newHandler(config)

	// Create a record
	record := slog.Record{
		Time:    time.Now(),
		Level:   slog.LevelInfo,
		Message: "test message",
	}

	// Create a context
	ctx := context.Background()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Number of goroutines to create
	numGoroutines := 100

	// Add the number of goroutines to the WaitGroup
	wg.Add(numGoroutines)

	// Create multiple goroutines to access the handler concurrently
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			// Call the Handle method, which should now be thread-safe
			err := handler.Handle(ctx, record)
			if err != nil {
				t.Errorf("Handler.Handle() error = %v", err)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
