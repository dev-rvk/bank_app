package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	InfoLevel       = "info"
	PanicLevel      = "panic"
	FatalLevel      = "fatal"
	ErrorLevel      = "error"
	WarnLevel       = "warn"
	DebugLevel      = "debug"
	TraceLevel      = "trace"
	BADKEY          = "!BADKEY"
)
const (
	timeFormat           = TimestampFormat
	defaultLogRetentions = 24
)

type Attribute struct {
	Key   string
	Value string
}

// Logger :
type Logger interface {
	Init()
	Panic(message string, args ...string)
	Fatal(message string, args ...string)
	Error(message string, err error, args ...string)
	Warn(message string, args ...string)
	Info(message string, args ...string)
	Debug(message string, args ...string)
}

// LoggingService :
type LoggingService struct {
	// Level          string
	context        context.Context //nolint:containedctx
	Filename       string
	file           *os.File
	RotateDuration time.Duration
}
type LogOption struct {
	NodeId    string
	ReplicaId int
	ShardId   int
	Module    string
}

type LoggingServiceOptions struct {
	FilenamePrefix    string
	RetentionDuration time.Duration
}

func NewLogger(module string, args ...LoggingServiceOptions) Logger {
	var rotateDuration time.Duration
	var fileName string

	for _, arg := range args {

		if arg.FilenamePrefix != "" {
			fileName = arg.FilenamePrefix
		}
		rotateDuration = defaultLogRetentions * time.Hour
		if arg.RetentionDuration > 0 {
			rotateDuration = arg.RetentionDuration
		}
		break
	}
	logService := &LoggingService{
		//  Level:          level,
		Filename:       fileName,
		RotateDuration: rotateDuration,
	}
	logService.Init()
	ctx, err := logService.getCtx()
	if err != nil {
		// slog.Info("fail to start Logger: %v", err)

		return nil
	}
	logService.context = ctx

	return logService
}

func SetLogLevel(level string) {
	switch level {
	case PanicLevel:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case FatalLevel:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case ErrorLevel:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case WarnLevel:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case InfoLevel:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case DebugLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case TraceLevel:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// Init :
func (loggingService *LoggingService) Init() {

	zerolog.TimeFieldFormat = timeFormat

	if loggingService.Filename != "" {
		if err := loggingService.openNew(); err != nil {

			return
		}
		if loggingService.RotateDuration != time.Duration(0) {
			runWatchdog(loggingService, loggingService.RotateDuration)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        loggingService.file,
			TimeFormat: time.RFC3339,
		})

	} else {
		fmt.Printf("\"file\": %v\n", "file")
		log.Logger = log.Output((zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: timeFormat,
		}))
		loggingService.file = os.Stdout
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: timeFormat,
		})
	}
}

// Panic :
func (loggingService *LoggingService) Panic(message string, args ...string) {

	// logString := log.Ctx(loggingService.context).Panic()
	logString := log.Panic().Ctx(loggingService.context)
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Fatal :
func (loggingService *LoggingService) Fatal(message string, args ...string) {

	logString := log.Ctx(loggingService.context).Fatal()
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)

}

// Error :
func (loggingService *LoggingService) Error(message string, err error, args ...string) {

	logString := log.Error().Ctx(loggingService.context).Caller(1).Stack().Err(err)
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Warn :
func (loggingService *LoggingService) Warn(message string, args ...string) {
	logString := log.Warn().Ctx(loggingService.context).Caller(1).Stack()
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Info :
func (loggingService *LoggingService) Info(message string, args ...string) {
	logString := log.Info().Ctx(loggingService.context).Caller(1).Stack()
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Debug :
func (loggingService *LoggingService) Debug(message string, args ...string) {

	logString := log.Debug().Ctx(loggingService.context).Caller(1).Stack()
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Rotate :
func (loggingService *LoggingService) Rotate() error {

	return loggingService.openNew()
}

// openNew :
func (loggingService *LoggingService) openNew() (err error) {
	if loggingService.Filename != "" {
		newName := nextName(loggingService.Filename)
		loggingService.file, err = os.OpenFile(newName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		return
	}
	loggingService.file = os.Stdout // default output

	return
}

// Close :
func (loggingService *LoggingService) Close() error {
	return loggingService.file.Close()

}

// getCtx :
func (loggingService *LoggingService) getCtx() (context.Context, error) {
	if _, err := os.Stat(loggingService.file.Name()); err != nil {
		err := loggingService.openNew()
		if err != nil {

			return nil, err
		}
	}

	return zerolog.New(loggingService.file).With().Timestamp().
		CallerWithSkipFrameCount(3).
		Logger().
		WithContext(context.Background()), nil
}

// nextName : Generic function to get filename with timestamp
func nextName(name string) string {
	filename := name
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]

	timestamp := time.Now().Format(TimestampFormat)

	return filepath.Join(fmt.Sprintf("%s-%s%s", prefix, timestamp, ext)) //nolint:gocritic
}

// watchdog : Structure which keeps track of duration after
type watchdog struct {
	interval time.Duration
}

// Run :
func (watchdog *watchdog) Run(l *LoggingService) {
	ticker := time.NewTicker(watchdog.interval)
	for range ticker.C {
		if err := l.Rotate(); err != nil {

			return
		}
	}
}

// runWatchdog :
func runWatchdog(loggingService *LoggingService, interval time.Duration) {
	wd := watchdog{interval: interval}
	go wd.Run(loggingService)
}

func argsToAttr(args []string) (Attribute, []string) {
	if len(args) == 1 {
		return Attribute{Key: BADKEY, Value: args[0]}, nil
	}
	return Attribute{Key: args[0], Value: args[1]}, args[2:]
}
