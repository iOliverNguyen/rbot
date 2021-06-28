package l

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/olvrng/rbot/be/pkg/dot"
)

// these constants are exported from zap
const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

type Field = zap.Field

var Any = zap.Any
var Array = zap.Array
var Binary = zap.Binary
var Bool = zap.Bool
var Boolp = zap.Boolp
var Bools = zap.Bools
var ByteString = zap.ByteString
var ByteStrings = zap.ByteStrings
var Complex128 = zap.Complex128
var Complex128p = zap.Complex128p
var Complex128s = zap.Complex128s
var Complex64 = zap.Complex64
var Complex64p = zap.Complex64p
var Complex64s = zap.Complex64s
var Duration = zap.Duration
var Durationp = zap.Durationp
var Durations = zap.Durations
var Error = zap.Error
var Errors = zap.Errors
var Float32 = zap.Float32
var Float32p = zap.Float32p
var Float32s = zap.Float32s
var Float64 = zap.Float64
var Float64p = zap.Float64p
var Float64s = zap.Float64s
var Inline = zap.Inline
var Int = zap.Int
var Int16 = zap.Int16
var Int16p = zap.Int16p
var Int16s = zap.Int16s
var Int32 = zap.Int32
var Int32p = zap.Int32p
var Int32s = zap.Int32s
var Int64 = zap.Int64
var Int64p = zap.Int64p
var Int64s = zap.Int64s
var Int8 = zap.Int8
var Int8p = zap.Int8p
var Int8s = zap.Int8s
var Intp = zap.Intp
var Ints = zap.Ints
var NamedError = zap.NamedError
var Namespace = zap.Namespace
var Object = zap.Object
var Reflect = zap.Reflect
var Skip = zap.Skip
var Stack = zap.Stack
var StackSkip = zap.StackSkip
var String = zap.String
var Stringer = zap.Stringer
var Stringp = zap.Stringp
var Strings = zap.Strings
var Time = zap.Time
var Timep = zap.Timep
var Times = zap.Times
var Uint = zap.Uint
var Uint16 = zap.Uint16
var Uint16p = zap.Uint16p
var Uint16s = zap.Uint16s
var Uint32 = zap.Uint32
var Uint32p = zap.Uint32p
var Uint32s = zap.Uint32s
var Uint64 = zap.Uint64
var Uint64p = zap.Uint64p
var Uint64s = zap.Uint64s
var Uint8 = zap.Uint8
var Uint8p = zap.Uint8p
var Uint8s = zap.Uint8s
var Uintp = zap.Uintp
var Uintptr = zap.Uintptr
var Uintptrp = zap.Uintptrp
var Uintptrs = zap.Uintptrs
var Uints = zap.Uints

func ID(key string, id dot.IntID) zap.Field {
	return zap.Int64(key, int64(id))
}
