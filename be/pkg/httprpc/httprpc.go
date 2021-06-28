package httprpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/olvrng/rbot/be/pkg/xerrors"
)

type Server interface {
	http.Handler

	PathPrefix() string

	// WithHooks returns a new server which calls given hooks on processing http
	// requests. It must always clone the server and leave the original one
	// unaffected.
	WithHooks(HooksBuilder) Server
}

type Muxer interface {
	Handle(path string, Handler http.Handler)
}

func BadRouteError(msg string, method, url string) *xerrors.APIError {
	return xerrors.Errorf(xerrors.BadRoute, nil, msg).
		WithMeta("invalid_route", method+" "+url)
}

func internalError(msg string) *xerrors.APIError {
	return xerrors.Errorf(xerrors.Internal, nil, msg)
}

func malformedRequestError(msg string) *xerrors.APIError {
	return xerrors.Errorf(xerrors.Malformed, nil, msg)
}

// JSON serialization for errors
type twerrJSON struct {
	Code string            "json:\"code\""
	Msg  string            "json:\"msg\""
	Meta map[string]string "json:\"meta,omitempty\""
}

// marshalErrorToJSON returns JSON from a xerrors.ErrorInterface, that can be used as HTTP error response body.
// If serialization fails, it will use a descriptive Internal error instead.
func marshalErrorToJSON(twerr ErrorInterface) []byte {
	// make sure that msg is not too large
	msg := twerr.Msg()
	if len(msg) > 1000000 {
		msg = msg[:1000000]
	}

	tj := twerrJSON{
		Code: twerr.Code().String(),
		Msg:  msg,
		Meta: twerr.MetaMap(),
	}

	buf, err := json.Marshal(&tj)
	if err != nil {
		buf = []byte("{\"type\": \"internal\", \"msg\": \"There was an error but it could not be serialized into JSON\"}") // fallback
	}
	return buf
}

func WriteError(ctx context.Context, resp http.ResponseWriter, hooks Hooks, info HookInfo, err error) {
	ctx, err = hooks.Error(ctx, info, err)

	twerr := TwirpError(err)
	statusCode := ServerHTTPStatusFromErrorCode(twerr.Code())

	respBody := marshalErrorToJSON(twerr)
	resp.Header().Set("Content-Type", "application/json") // Error responses are always JSON
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBody)))
	resp.WriteHeader(statusCode) // set HTTP status code and send response

	_, writeErr := resp.Write(respBody)
	if writeErr != nil {
		// We have three options here. We could log the error, call the Error
		// hook, or just silently ignore the error.
		//
		// Logging is unacceptable because we don't have a user-controlled
		// logger; writing out to stderr without permission is too rude.
		//
		// Calling the Error hook would confuse users: it would mean the Error
		// hook got called twice for one request, which is likely to lead to
		// duplicated log messages and metrics, no matter how well we document
		// the behavior.
		//
		// Silently ignoring the error is our least-bad option. It's highly
		// likely that the connection is broken and the original 'err' says
		// so anyway.
		_ = writeErr
	}
}

// internalWithCause is a Twirp Internal error wrapping an original error cause, accessible
// by github.com/pkg/errors.Cause, but the original error message is not exposed on Msg().
type internalWithCause struct {
	msg   string
	cause error
}

func (e *internalWithCause) Cause() error                                   { return e.cause }
func (e *internalWithCause) Error() string                                  { return e.msg + ": " + e.cause.Error() }
func (e *internalWithCause) Code() xerrors.Code                             { return xerrors.Internal }
func (e *internalWithCause) Msg() string                                    { return e.msg }
func (e *internalWithCause) Meta(key string) string                         { return "" }
func (e *internalWithCause) MetaMap() map[string]string                     { return nil }
func (e *internalWithCause) WithMeta(key string, val string) ErrorInterface { return e }

// ensurePanicResponses makes sure that rpc methods causing a panic still result in a Twirp Internal
// error response (status 500), and error hooks are properly called with the panic wrapped as an error.
// The panic is re-raised so it can be handled normally with middleware.
func ensurePanicResponses(ctx context.Context, resp http.ResponseWriter, hooks Hooks, info HookInfo) {
	if r := recover(); r != nil {
		// Wrap the panic as an error so it can be passed to error hooks.
		// The original error is accessible from error hooks, but not visible in the response.
		err := errFromPanic(r)
		twerr := &internalWithCause{msg: "Internal service panic", cause: err}
		// Actually write the error
		WriteError(ctx, resp, hooks, info, twerr)
		// If possible, flush the error to the wire.
		f, ok := resp.(http.Flusher)
		if ok {
			f.Flush()
		}
		panic(r)
	}
}

// errFromPanic returns the typed error if the recovered panic is an error, otherwise formats as error.
func errFromPanic(p interface{}) error {
	if err, ok := p.(error); ok {
		return err
	}
	return fmt.Errorf("panic: %v", p)
}

type ExecFunc func(context.Context) (ctx context.Context, respContent Message, err error)
type ServeFunc func(ctx context.Context, resp http.ResponseWriter, req *http.Request, hooks Hooks, info *HookInfo, reqContent Message, fn ExecFunc)

func ParseRequestHeader(req *http.Request) (ServeFunc, error) {
	if req.Method != "POST" {
		msg := fmt.Sprintf("unsupported method %q (only POST is allowed)", req.Method)
		return nil, BadRouteError(msg, req.Method, req.URL.Path)
	}
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		return ServeJSON, nil
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		return nil, BadRouteError(msg, req.Method, req.URL.Path)
	}
}

func ServeJSON(
	ctx context.Context,
	resp http.ResponseWriter,
	req *http.Request,
	hooks Hooks,
	info *HookInfo,
	reqContent Message,
	fn ExecFunc,
) {
	if err := json.NewDecoder(req.Body).Decode(reqContent); err != nil {
		WriteError(ctx, resp, hooks, *info, malformedRequestError("the json request could not be decoded").WithMeta("cause", err.Error()))
		return
	}
	info.Request = reqContent

	var err error
	var respContent Message
	func() {
		defer ensurePanicResponses(ctx, resp, hooks, *info)
		ctx, respContent, err = fn(ctx)
	}()
	if err != nil {
		WriteError(ctx, resp, hooks, *info, err)
		return
	}
	if respContent == nil {
		WriteError(ctx, resp, hooks, *info, internalError("received a nil response"))
		return
	}
	info.Response = respContent
	ctx, err = hooks.ResponsePrepared(ctx, *info, resp.Header())
	if err != nil {
		WriteError(ctx, resp, hooks, *info, err)
		return
	}

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(respContent); err != nil {
		WriteError(ctx, resp, hooks, *info, xerrors.Errorf(xerrors.Internal, err, "failed to marshal json response: %v"))
		return
	}
	respBytes := buf.Bytes()
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)
	defer hooks.ResponseSent(ctx, *info)
	if _, err = resp.Write(respBytes); err != nil {
		return
	}
}
