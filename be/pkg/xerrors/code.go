package xerrors

type Code int

func (c Code) String() string {
	code := mapCodes[c]
	if code == "" {
		return "unknown"
	}
	return code
}

// Error constants from twirp
const (
	NoError            = Code(0)
	Canceled           = Code(1)
	Unknown            = Code(2)
	InvalidArgument    = Code(3)
	DeadlineExceeded   = Code(4)
	NotFound           = Code(5)
	AlreadyExists      = Code(6)
	PermissionDenied   = Code(7)
	ResourceExhausted  = Code(8)
	FailedPrecondition = Code(9)
	Aborted            = Code(10)
	OutOfRange         = Code(11)
	Unimplemented      = Code(12)
	Internal           = Code(13)
	Unavailable        = Code(14)
	DataLoss           = Code(15)
	Unauthenticated    = Code(16)

	RuntimePanic = Code(100)
	BadRoute     = Code(101)
	Malformed    = Code(102)
)

var (
	mapCodes [Unauthenticated + 1]string
)

func init() {
	mapCodes[Canceled] = "canceled"
	mapCodes[Unknown] = "unknown"
	mapCodes[InvalidArgument] = "invalid_argument"
	mapCodes[DeadlineExceeded] = "deadline_exceeded"
	mapCodes[NotFound] = "not_found"
	mapCodes[AlreadyExists] = "already_exists"
	mapCodes[PermissionDenied] = "permission_denied"
	mapCodes[Unauthenticated] = "unauthenticated"
	mapCodes[ResourceExhausted] = "resource_exhausted"
	mapCodes[FailedPrecondition] = "failed_precondition"
	mapCodes[Aborted] = "aborted"
	mapCodes[OutOfRange] = "out_of_range"
	mapCodes[Unimplemented] = "unimplemented"
	mapCodes[Internal] = "internal"
	mapCodes[Unavailable] = "unavailable"
	mapCodes[DataLoss] = "data_loss"
	mapCodes[NoError] = "ok"
}

func IsValidStandardErrorCode(c Code) bool {
	return c >= 0 && int(c) < len(mapCodes)
}

func DefaultErrorMessage(code Code) string {
	switch code {
	case NoError:
		return ""
	case NotFound:
		return "Not Found"
	case InvalidArgument:
		return "Invalid Argument"
	case Internal:
		return "Internal Error"
	case Unauthenticated:
		return "Unauthenticated. Please login."
	case PermissionDenied:
		return "Permission Denied. Please check your permission."
	case Unimplemented:
		return "The request is unimplemented yet."
	}
	return "Internal Error"
}
