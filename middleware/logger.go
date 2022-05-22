package middleware

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LoggerOption = func(mw *Middleware)

func WithLogger(logger *zap.Logger) LoggerOption {
	return func(mw *Middleware) {
		mw.logger = logger
	}
}

func WithCustomFields(fieldHandler func(c *fiber.Ctx, start time.Time, stop time.Time) []zap.Field) LoggerOption {
	return func(mw *Middleware) {
		mw.fieldHandler = fieldHandler
	}
}

func WithPanicHandler(stackTraceHandler func(c *fiber.Ctx, e interface{})) LoggerOption {
	return func(mw *Middleware) {
		mw.stackTraceHandler = stackTraceHandler
	}
}

type Middleware struct {
	logger            *zap.Logger
	fieldHandler      func(c *fiber.Ctx, start time.Time, stop time.Time) []zap.Field
	stackTraceHandler func(c *fiber.Ctx, e interface{})
}

func NewLogger(opts ...LoggerOption) *Middleware {
	mw := &Middleware{}

	for _, opt := range opts {
		opt(mw)
	}

	if mw.logger == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("failed to load zap logger: %v", err)
		}
		mw.logger = logger
	}

	if mw.fieldHandler == nil {
		// zap logger handler for access logging
		mw.fieldHandler = func(c *fiber.Ctx, start time.Time, stop time.Time) []zap.Field {
			return []zap.Field{
				zap.String("method", string(c.Request().Header.Method())),
				zap.String("url", string(c.Request().URI().FullURI())),
				zap.Int("status", c.Response().StatusCode()),
				zap.Int("bytes", len(c.Response().Body())),
				zap.Int64("elapsed", stop.Sub(start).Microseconds()),
			}
		}
	}

	if mw.stackTraceHandler == nil {
		// zap panic logger handler
		mw.stackTraceHandler = func(c *fiber.Ctx, e interface{}) {

			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, false)]
			mw.logger.Error(
				"panic",
				zap.String("method", string(c.Request().Header.Method())),
				zap.String("url", string(c.Request().URI().FullURI())),
				zap.Int("status", 500),
				zap.Int("bytes", 0),
				zap.String("stack", fmt.Sprintf("panic: %v\n%s\n", e, buf)),
			)
		}
	}

	return mw
}

func (mw *Middleware) Logger() *zap.Logger {
	return mw.logger
}

func (mw *Middleware) StackTraceHandler() func(c *fiber.Ctx, e interface{}) {
	return mw.stackTraceHandler
}

func (mw *Middleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		start := time.Now()

		chainErr := c.Next()
		if chainErr != nil {
			if err := c.App().ErrorHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		stop := time.Now()

		if c.Response().StatusCode() < 400 {
			mw.logger.Info(
				"http request",
				mw.fieldHandler(c, start, stop)...,
			)
		} else if c.Response().StatusCode() < 500 {
			mw.logger.Warn(
				"http request",
				mw.fieldHandler(c, start, stop)...,
			)
		} else {
			mw.logger.Error(
				"http request",
				mw.fieldHandler(c, start, stop)...,
			)
		}

		return nil
	}
}
