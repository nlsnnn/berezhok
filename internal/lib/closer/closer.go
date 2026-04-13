package closer

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

// closeFn - структура для хранения функции закрытия и её имени
type closeFn struct {
	name string
	fn   func(context.Context) error
}

// closer — структура для управления функциями закрытия ресурсов.
type closer struct {
	mu    sync.Mutex // защищает слайс funcs от конкурентной записи
	once  sync.Once  // гарантирует что CloseAll выполнится только один раз
	funcs []closeFn  // накопленные функции закрытия в порядке добавления
	log   *slog.Logger
}

// globalCloser - глобальный экземпляр closer-а для всего приложения.
var globalCloser = &closer{}

func SetLogger(log *slog.Logger) {
	globalCloser.log = log
}

// Add добавляет функцию закрытия в глобальный closer.
func Add(name string, fn func(context.Context) error) {
	globalCloser.add(name, fn)
}

// CloseAll вызывает все функции закрытия глобального closer-а в обратном порядке (LIFO).
func CloseAll(ctx context.Context) error {
	return globalCloser.closeAll(ctx)
}

// add добавляет функцию закрытия с именем ресурса.
func (c *closer) add(name string, fn func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, closeFn{name: name, fn: fn})
}

// closeAll вызывает все зарегистрированные функции закрытия в обратном порядке (LIFO).
func (c *closer) closeAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			return
		}

		c.log.Info("start shutdown", "count", len(funcs))

		var errs []error

		// Идём от конца к началу — LIFO, как defer.
		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]

			start := time.Now()
			c.log.Info("closing resource", "name", f.name)

			if err := f.fn(ctx); err != nil {
				// Логируем ошибку, но продолжаем закрывать остальные ресурсы.
				c.log.Error("error while closing resource",
					"name", f.name,
					"error", err,
					"duration", time.Since(start),
				)

				errs = append(errs, err)
			} else {
				c.log.Info("resource closed", "name", f.name, "duration", time.Since(start))
			}
		}

		c.log.Info("all resources closed")
		result = errors.Join(errs...)
	})

	return result
}
