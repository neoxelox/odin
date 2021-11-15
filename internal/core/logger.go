package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	gommon "github.com/labstack/gommon/log"
	"github.com/neoxelox/odin/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

var (
	ZlevelToGlevel = map[zerolog.Level]gommon.Lvl{
		zerolog.DebugLevel: gommon.DEBUG,
		zerolog.InfoLevel:  gommon.INFO,
		zerolog.WarnLevel:  gommon.WARN,
		zerolog.ErrorLevel: gommon.ERROR,
		zerolog.Disabled:   gommon.OFF,
	}

	GlevelToZlevel = map[gommon.Lvl]zerolog.Level{
		gommon.DEBUG: zerolog.DebugLevel,
		gommon.INFO:  zerolog.InfoLevel,
		gommon.WARN:  zerolog.WarnLevel,
		gommon.ERROR: zerolog.ErrorLevel,
		gommon.OFF:   zerolog.Disabled,
	}

	ZlevelToPlevel = map[zerolog.Level]pgx.LogLevel{
		zerolog.TraceLevel: pgx.LogLevelTrace,
		zerolog.DebugLevel: pgx.LogLevelDebug,
		zerolog.InfoLevel:  pgx.LogLevelInfo,
		zerolog.WarnLevel:  pgx.LogLevelWarn,
		zerolog.ErrorLevel: pgx.LogLevelError,
		zerolog.Disabled:   pgx.LogLevelNone,
	}

	PlevelToZlevel = map[pgx.LogLevel]zerolog.Level{
		pgx.LogLevelTrace: zerolog.TraceLevel,
		pgx.LogLevelDebug: zerolog.DebugLevel,
		pgx.LogLevelInfo:  zerolog.InfoLevel,
		pgx.LogLevelWarn:  zerolog.WarnLevel,
		pgx.LogLevelError: zerolog.ErrorLevel,
		pgx.LogLevelNone:  zerolog.Disabled,
	}
)

type Logger struct {
	configuration internal.Configuration
	logger        zerolog.Logger
	level         zerolog.Level
	out           io.Writer
	prefix        string
	header        string
	verbose       bool
}

func NewLogger(configuration internal.Configuration) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "timestamp"
	zerolog.CallerSkipFrameCount = 3

	level := zerolog.DebugLevel
	if configuration.Environment != internal.Environment.DEVELOPMENT {
		level = zerolog.InfoLevel
	}

	out := diode.NewWriter(os.Stderr, 1000, 10*time.Millisecond, func(missed int) {
		fmt.Fprintf(os.Stderr, "Logger dropped %d messages", missed)
	})

	return &Logger{
		configuration: configuration,
		logger:        zerolog.New(out).With().Str("service", configuration.AppName).Timestamp().Logger().Level(level),
		level:         level,
		out:           out,
		prefix:        configuration.AppName,
		header:        "",
		verbose:       level == zerolog.DebugLevel,
	}
}

func (self Logger) Logger() *zerolog.Logger {
	return &self.logger
}

func (self *Logger) SetLogger(l zerolog.Logger) {
	self.logger = l
}

func (self Logger) Flush() {
	os.Stderr.Sync()
}

func (self Logger) Close(ctx context.Context) error {
	self.Flush()
	dw, _ := self.out.(diode.Writer)
	return dw.Close()
}

func (self Logger) Output() io.Writer {
	return self.out
}

func (self *Logger) SetOutput(w io.Writer) {
	self.logger = self.logger.Output(w)
	self.out = w
}

func (self Logger) Prefix() string {
	return self.prefix
}

func (self *Logger) SetPrefix(p string) {
	self.prefix = p
}

func (self Logger) GLevel() gommon.Lvl {
	return ZlevelToGlevel[self.level]
}

func (self *Logger) SetGLevel(l gommon.Lvl) {
	zlevel := GlevelToZlevel[l]
	self.logger = self.logger.Level(zlevel)
	self.level = zlevel
}

func (self Logger) PLevel() pgx.LogLevel {
	return ZlevelToPlevel[self.level]
}

func (self *Logger) SetPLevel(l pgx.LogLevel) {
	zlevel := PlevelToZlevel[l]
	self.logger = self.logger.Level(zlevel)
	self.level = zlevel
}

func (self Logger) ZLevel() zerolog.Level {
	return self.level
}

func (self *Logger) SetZLevel(l zerolog.Level) {
	zlevel := l
	self.logger = self.logger.Level(zlevel)
	self.level = zlevel
}

func (self *Logger) Header() string {
	return self.header
}

func (self *Logger) SetHeader(h string) {
	self.header = h
}

func (self Logger) Verbose() bool {
	return self.verbose
}

func (self Logger) SetVerbose(v bool) {
	self.verbose = v
}

func (self Logger) Print(i ...interface{}) {
	self.logger.Log().Msg(fmt.Sprint(i...))
}

func (self Logger) Printf(format string, i ...interface{}) {
	self.logger.Log().Msgf(format, i...)
}

func (self Logger) Debug(i ...interface{}) {
	self.logger.Debug().Msg(fmt.Sprint(i...))
}

func (self Logger) Debugf(format string, i ...interface{}) {
	self.logger.Debug().Msgf(format, i...)
}

func (self Logger) Info(i ...interface{}) {
	self.logger.Info().Msg(fmt.Sprint(i...))
}

func (self Logger) Infof(format string, i ...interface{}) {
	self.logger.Info().Msgf(format, i...)
}

func (self Logger) Warn(i ...interface{}) {
	self.logger.Warn().Msg(fmt.Sprint(i...))
}

func (self Logger) Warnf(format string, i ...interface{}) {
	self.logger.Warn().Msgf(format, i...)
}

func (self Logger) Error(i ...interface{}) {
	if self.configuration.Environment == internal.Environment.PRODUCTION {
		self.logger.Error().Msg(fmt.Sprint(i...))
	} else {
		if i != nil {
			err, ok := i[0].(*internal.Error)
			if !ok {
				fmt.Printf("%+v\n", i...)
			} else {
				fmt.Printf("%+v\n", err.Outer())
			}
		}
	}
}

func (self Logger) Errorf(format string, i ...interface{}) {
	self.logger.Error().Msgf(format, i...)
}

func (self Logger) Fatal(i ...interface{}) {
	self.logger.Fatal().Msg(fmt.Sprint(i...))
}

func (self Logger) Fatalf(format string, i ...interface{}) {
	self.logger.Fatal().Msgf(format, i...)
}

func (self Logger) Panic(i ...interface{}) {
	self.logger.Panic().Msg(fmt.Sprint(i...))
}

func (self Logger) Panicf(format string, i ...interface{}) {
	self.logger.Panic().Msgf(format, i...)
}
