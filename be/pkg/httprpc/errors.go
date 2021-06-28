package httprpc

import "github.com/olvrng/rbot/be/pkg/xerrors"

func ServerHTTPStatusFromErrorCode(code xerrors.Code) int {
	switch code {
	case xerrors.Canceled:
		return 408 // RequestTimeout
	case xerrors.Unknown:
		return 500 // Internal Server Error
	case xerrors.InvalidArgument:
		return 400 // BadRequest
	case xerrors.Malformed:
		return 400 // BadRequest
	case xerrors.DeadlineExceeded:
		return 408 // RequestTimeout
	case xerrors.NotFound:
		return 404 // Not Found
	case xerrors.BadRoute:
		return 404 // Not Found
	case xerrors.AlreadyExists:
		return 409 // Conflict
	case xerrors.PermissionDenied:
		return 403 // Forbidden
	case xerrors.Unauthenticated:
		return 401 // Unauthorized
	case xerrors.ResourceExhausted:
		return 403 // Forbidden
	case xerrors.FailedPrecondition:
		return 412 // Precondition Failed
	case xerrors.Aborted:
		return 409 // Conflict
	case xerrors.OutOfRange:
		return 400 // Bad Request
	case xerrors.Unimplemented:
		return 501 // Not Implemented
	case xerrors.Internal:
		return 500 // Internal Server Error
	case xerrors.Unavailable:
		return 503 // Service Unavailable
	case xerrors.DataLoss:
		return 500 // Internal Server Error
	case xerrors.NoError:
		return 200 // OK
	default:
		return 0 // Invalid!
	}
}

func TwirpError(err error) TwError {
	if err == nil {
		return nil
	}
	if xerr, ok := err.(TwError); ok {
		return xerr
	}
	xerr, ok := err.(*xerrors.APIError)
	if !ok {
		xerr = xerrors.Errorf(xerrors.Internal, err, "")
	}

	// FIXED: fatal error: concurrent map iteration and map write
	//
	// This clones the meta map from the inner error. We don't reuse the map to
	// prevent concurrency access and make the function TwirpError thread-safe.
	//
	// TODO(vu): Remove the TwError and convert to ErrorJSON directly.
	meta := map[string]string{}
	for k, v := range xerr.Meta {
		meta[k] = v
	}
	if xerr.Err != nil {
		meta["cause"] = xerr.Err.Error()
	}
	if xerr.Original != "" {
		meta["orig"] = xerr.Original
	}
	return twError{xerr, meta}
}
