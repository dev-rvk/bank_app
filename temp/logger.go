package ports

import (
	"context"
	"os"
	"time"
)

type Logger interface {
	Init()
	Panic(message string, attributes ...any)
	Fatal(message string, attributes ...any)
	Error(message string, err error, attributes ...any)
	Warn(message string, attributes ...any)
	WarnErr(message string, err error, attributes ...any)
	Info(message string, attributes ...any)
	Debug(message string, attributes ...any)
	DebugErr(message string, err error, attributes ...any)
	SetCallerSkip(skip int)
}

type LoggingService struct {
	context        context.Context
	file           *os.File
	Filename       string
	RotateDuration time.Duration
	CallerSkip     int
}

type LogOption struct {
	NodeID    string
	Module    string
	ReplicaID int
	ShardID   int
}

type LoggingServiceOptions struct {
	FilenamePrefix    string
	RetentionDuration time.Duration
}

type Attribute struct {
	Value any
	Key   string
}
