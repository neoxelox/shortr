package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gommon "github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

var (
	zlevelToGlevel = map[zerolog.Level]gommon.Lvl{
		zerolog.DebugLevel: gommon.DEBUG,
		zerolog.InfoLevel:  gommon.INFO,
		zerolog.WarnLevel:  gommon.WARN,
		zerolog.ErrorLevel: gommon.ERROR,
		zerolog.Disabled:   gommon.OFF,
	}

	glevelToZlevel = map[gommon.Lvl]zerolog.Level{
		gommon.DEBUG: zerolog.DebugLevel,
		gommon.INFO:  zerolog.InfoLevel,
		gommon.WARN:  zerolog.WarnLevel,
		gommon.ERROR: zerolog.ErrorLevel,
		gommon.OFF:   zerolog.Disabled,
	}

	zlevelToPlevel = map[zerolog.Level]pgx.LogLevel{
		zerolog.TraceLevel: pgx.LogLevelTrace,
		zerolog.DebugLevel: pgx.LogLevelDebug,
		zerolog.InfoLevel:  pgx.LogLevelInfo,
		zerolog.WarnLevel:  pgx.LogLevelWarn,
		zerolog.ErrorLevel: pgx.LogLevelError,
		zerolog.Disabled:   pgx.LogLevelNone,
	}

	plevelToZlevel = map[pgx.LogLevel]zerolog.Level{
		pgx.LogLevelTrace: zerolog.TraceLevel,
		pgx.LogLevelDebug: zerolog.DebugLevel,
		pgx.LogLevelInfo:  zerolog.InfoLevel,
		pgx.LogLevelWarn:  zerolog.WarnLevel,
		pgx.LogLevelError: zerolog.ErrorLevel,
		pgx.LogLevelNone:  zerolog.Disabled,
	}
)

// Logger implements the echo.Logger interface
type Logger struct {
	logger zerolog.Logger
	level  zerolog.Level
	out    io.Writer
	prefix string
}

// New creates a new Logger instance
func New(prefix string) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "timestamp"
	zerolog.CallerSkipFrameCount = 3

	out := diode.NewWriter(os.Stderr, 1000, 10*time.Millisecond, func(missed int) {
		fmt.Fprintf(os.Stderr, "Logger dropped %d messages", missed)
	})

	return &Logger{
		logger: zerolog.New(out).With().Str("prefix", prefix).Timestamp().Logger().Level(zerolog.InfoLevel),
		level:  zerolog.InfoLevel,
		out:    out,
		prefix: prefix,
	}
}

// Logger returns a copy of the internal logger
func (l Logger) Logger() zerolog.Logger {
	return l.logger
}

// Log satisfies the pgx.Logger interface
func (l Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	l.logger = l.logger.With().Fields(data).Logger()
	l.logger.WithLevel(plevelToZlevel[level]).Msg(msg)
}

// Output satisfies the echo.Logger interface
func (l Logger) Output() io.Writer {
	return l.out
}

// SetOutput satisfies the echo.Logger interface
func (l *Logger) SetOutput(w io.Writer) {
	l.logger = l.logger.Output(w)
	l.out = w
}

// Prefix satisfies the echo.Logger interface
func (l Logger) Prefix() string {
	return l.prefix
}

// SetPrefix satisfies the echo.Logger interface
func (l *Logger) SetPrefix(p string) {
	// Have to create a new logger, becayse zerolog doesn't dedup fields.
	// Otherwise "prefix" would appear twice in the log output.
	ll := New(p)
	l.logger = ll.logger
	l.level = ll.level
	l.out = ll.out
	l.prefix = ll.prefix
}

// Level satisfies the echo.Logger interface
func (l Logger) Level() gommon.Lvl {
	return zlevelToGlevel[l.level]
}

// SetLevel satisfies the echo.Logger interface
func (l *Logger) SetLevel(v gommon.Lvl) {
	zlevel := glevelToZlevel[v]
	l.logger = l.logger.Level(zlevel)
	l.level = zlevel
}

// SetHeader satisfies the echo.Logger interface
func (l *Logger) SetHeader(h string) {}

// Print satisfies the echo.Logger interface
func (l Logger) Print(i ...interface{}) {
	l.logger.Log().Msg(fmt.Sprint(i...))
}

// Printf satisfies the echo.Logger interface
func (l Logger) Printf(format string, i ...interface{}) {
	l.logger.Log().Msgf(format, i...)
}

// Printj satisfies the echo.Logger interface
func (l Logger) Printj(j gommon.JSON) {
	for k, v := range j {
		j, _ := json.Marshal(v)
		l.logger = l.logger.With().RawJSON(k, j).Logger()
	}
	l.logger.Log().Msg("")
}

// Debug satisfies the echo.Logger interface
func (l Logger) Debug(i ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprint(i...))
}

// Debugf satisfies the echo.Logger interface
func (l Logger) Debugf(format string, i ...interface{}) {
	l.logger.Debug().Msgf(format, i...)
}

// Debugj satisfies the echo.Logger interface
func (l Logger) Debugj(j gommon.JSON) {
	for k, v := range j {
		j, _ := json.Marshal(v)
		l.logger = l.logger.With().RawJSON(k, j).Logger()
	}
	l.logger.Debug().Msg("")
}

// Info satisfies the echo.Logger interface
func (l Logger) Info(i ...interface{}) {
	l.logger.Info().Msg(fmt.Sprint(i...))
}

// Infof satisfies the echo.Logger interface
func (l Logger) Infof(format string, i ...interface{}) {
	l.logger.Info().Msgf(format, i...)
}

// Infoj satisfies the echo.Logger interface
func (l Logger) Infoj(j gommon.JSON) {
	for k, v := range j {
		j, _ := json.Marshal(v)
		l.logger = l.logger.With().RawJSON(k, j).Logger()
	}
	l.logger.Info().Msg("")
}

// Warn satisfies the echo.Logger interface
func (l Logger) Warn(i ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprint(i...))
}

//Warnf satisfies the echo.Logger interface
func (l Logger) Warnf(format string, i ...interface{}) {
	l.logger.Warn().Msgf(format, i...)
}

// Warnj satisfies the echo.Logger interface
func (l Logger) Warnj(j gommon.JSON) {
	for k, v := range j {
		j, _ := json.Marshal(v)
		l.logger = l.logger.With().RawJSON(k, j).Logger()
	}
	l.logger.Warn().Msg("")
}

// Error satisfies the echo.Logger interface
func (l Logger) Error(i ...interface{}) {
	l.logger.Error().Msg(fmt.Sprint(i...))
}

// Errorf satisfies the echo.Logger interface
func (l Logger) Errorf(format string, i ...interface{}) {
	l.logger.Error().Msgf(format, i...)
}

// Errorj satisfies the echo.Logger interface
func (l Logger) Errorj(j gommon.JSON) {
	for k, v := range j {
		j, _ := json.Marshal(v)
		l.logger = l.logger.With().RawJSON(k, j).Logger()
	}
	l.logger.Error().Msg("")
}

// Fatal satisfies the echo.Logger interface
func (l Logger) Fatal(i ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprint(i...))
}

// Fatalf satisfies the echo.Logger interface
func (l Logger) Fatalf(format string, i ...interface{}) {
	l.logger.Fatal().Msgf(format, i...)
}

// Fatalj satisfies the echo.Logger interface
func (l Logger) Fatalj(j gommon.JSON) {
	for k, v := range j {
		j, _ := json.Marshal(v)
		l.logger = l.logger.With().RawJSON(k, j).Logger()
	}
	l.logger.Fatal().Msg("")
}

// Panic satisfies the echo.Logger interface
func (l Logger) Panic(i ...interface{}) {
	l.logger.Panic().Msg(fmt.Sprint(i...))
}

// Panicf satisfies the echo.Logger interface
func (l Logger) Panicf(format string, i ...interface{}) {
	l.logger.Panic().Msgf(format, i...)
}

// Panicj satisfies the echo.Logger interface
func (l Logger) Panicj(j gommon.JSON) {
	for k, v := range j {
		j, _ := json.Marshal(v)
		l.logger = l.logger.With().RawJSON(k, j).Logger()
	}
	l.logger.Panic().Msg("")
}

// Standard implements echo.Logger interface
func Standard(logger *Logger) *Logger {
	return &Logger{
		logger: logger.logger.With().Caller().Logger(),
		level:  zerolog.InfoLevel,
		out:    logger.out,
		prefix: logger.prefix,
	}
}

// Database implements pgx.Logger interface
func Database(logger *Logger) *Logger {
	return &Logger{
		logger: logger.logger.With().Str("module", "pgx").Logger(),
		level:  zerolog.ErrorLevel,
		out:    logger.out,
		prefix: logger.prefix,
	}
}

// Middleware implements echo.MiddlewareFunc interface
func Middleware(logger *Logger) echo.MiddlewareFunc {
	l := logger.Logger()
	skipper := middleware.DefaultSkipper
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if skipper(ctx) {
				return next(ctx)
			}

			req := ctx.Request()
			res := ctx.Response()

			start := time.Now()

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			stop := time.Now()

			l.Info().
				Str("method", req.Method).
				Str("path", req.RequestURI).
				Int("status", res.Status).
				Str("ip_address", ctx.RealIP()).
				Str("user_agent", req.UserAgent()).
				Dur("latency", stop.Sub(start)).
				Msg("")

			return nil
		}
	}
}
