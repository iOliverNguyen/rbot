package httprpc

import (
	"context"
	"net/http"
)

type HookInfo struct {
	Route       string
	HTTPRequest *http.Request
	Request     Message
	Response    Message
	Inner       interface{}
}

type Hooks struct {
	RequestReceived  func(ctx context.Context, info HookInfo) (context.Context, error)
	RequestRouted    func(ctx context.Context, info HookInfo) (context.Context, error)
	ResponsePrepared func(ctx context.Context, info HookInfo, respHeaders http.Header) (context.Context, error)
	ResponseSent     func(ctx context.Context, info HookInfo)
	Error            func(ctx context.Context, info HookInfo, err error) (context.Context, error)
}

type HooksBuilder interface {
	BuildHooks() Hooks
}

type HooksFunc func() Hooks

func (h HooksFunc) BuildHooks() Hooks { return h() }

type chainHooks []HooksBuilder

func (s chainHooks) BuildHooks() Hooks {
	switch len(s) {
	case 0:
		return Hooks{}
	case 1:
		return s[0].BuildHooks()
	}
	hooks := make([]Hooks, len(s))
	for i, b := range s {
		hooks[i] = b.BuildHooks()
	}
	return Hooks{
		RequestReceived: func(ctx context.Context, info HookInfo) (_ context.Context, err error) {
			for _, h := range hooks {
				if h.RequestReceived != nil {
					ctx, err = h.RequestReceived(ctx, info)
					if err != nil {
						return ctx, err
					}
				}
			}
			return ctx, nil
		},
		RequestRouted: func(ctx context.Context, info HookInfo) (_ context.Context, err error) {
			for _, h := range hooks {
				if h.RequestRouted != nil {
					ctx, err = h.RequestRouted(ctx, info)
					if err != nil {
						return ctx, err
					}
				}
			}
			return ctx, nil
		},
		ResponsePrepared: func(ctx context.Context, info HookInfo, respHeaders http.Header) (_ context.Context, err error) {
			for _, h := range hooks {
				if h.ResponsePrepared != nil {
					ctx, err = h.ResponsePrepared(ctx, info, respHeaders)
					if err != nil {
						return ctx, err
					}
				}
			}
			return ctx, nil
		},
		ResponseSent: func(ctx context.Context, info HookInfo) {
			for _, h := range hooks {
				if h.ResponseSent != nil {
					h.ResponseSent(ctx, info)
				}
			}
		},
		Error: func(ctx context.Context, info HookInfo, err error) (context.Context, error) {
			for _, h := range hooks {
				if h.Error != nil {
					ctx, err = h.Error(ctx, info, err)
				}
			}
			return ctx, err
		},
	}
}

func ChainHooks(hooks ...HooksBuilder) HooksBuilder {
	length := 2 * len(hooks)
	if length < 8 {
		length = 8
	}
	res := make(chainHooks, 0, length)
	for _, h := range hooks {
		if h == nil {
			continue
		}
		if hs, ok := h.(chainHooks); ok {
			res = append(res, hs...)
		} else {
			res = append(res, h)
		}
	}
	switch len(res) {
	case 0:
		return chainHooks{}
	case 1:
		return res[0]
	default:
		return res
	}
}

func WrapHooks(builder HooksBuilder) (res Hooks) {
	var hooks Hooks
	if builder != nil {
		hooks = builder.BuildHooks()
	}
	if hooks.RequestReceived == nil {
		hooks.RequestReceived = func(ctx context.Context, _ HookInfo) (context.Context, error) { return ctx, nil }
	}
	if hooks.RequestRouted == nil {
		hooks.RequestRouted = func(ctx context.Context, _ HookInfo) (context.Context, error) { return ctx, nil }
	}
	if hooks.ResponsePrepared == nil {
		hooks.ResponsePrepared = func(ctx context.Context, _ HookInfo, _ http.Header) (context.Context, error) { return ctx, nil }
	}
	if hooks.ResponseSent == nil {
		hooks.ResponseSent = func(ctx context.Context, _ HookInfo) {}
	}
	if hooks.Error == nil {
		hooks.Error = func(ctx context.Context, _ HookInfo, err error) (context.Context, error) { return ctx, err }
	}
	return hooks
}

func WithHooks(servers []Server, hooks ...HooksBuilder) []Server {
	if len(hooks) == 0 {
		return servers
	}
	result := make([]Server, len(servers))
	for i, s := range servers {
		result[i] = s.WithHooks(ChainHooks(hooks...))
	}
	return result
}
